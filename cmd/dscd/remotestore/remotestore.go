// Package remotestore serves module KV stores read-only from a remote (archive)
// node at a fixed height, through Tendermint abci_query "/store/<name>/key" and
// "/store/<name>/subspace". It lets offline tooling run keeper code against any
// historical state the archive node still has. Writes stay in memory.
//
// WARNING: in SDK v0.46 the IAVL "/subspace" query iterates the LATEST version
// whatever height is requested (store/iavl/store.go), so Store iterators and
// LoadPrefix* return current data. Only Client.Key / Client.KeysBatch are
// historical reads.
package remotestore

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/store/cachekv"
	"github.com/cosmos/cosmos-sdk/store/tracekv"
	"github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/types/kv"
	dbm "github.com/tendermint/tm-db"
)

// Client queries one archive node at one height.
type Client struct {
	RPC    string // e.g. https://archive-node.decimalchain.com/rpc
	Height int64
	HTTP   *http.Client
	mu     sync.Mutex
	Calls  int
	// Progress, if set, gets one line per remote query.
	Progress io.Writer
}

type abciResp struct {
	Result struct {
		Response struct {
			Code  int    `json:"code"`
			Log   string `json:"log"`
			Value []byte `json:"value"`
		} `json:"response"`
	} `json:"result"`
	Error *struct {
		Message string `json:"message"`
		Data    string `json:"data"`
	} `json:"error"`
}

func (c *Client) query(path string, data []byte) ([]byte, error) {
	q := url.Values{}
	q.Set("path", fmt.Sprintf("%q", path))
	q.Set("data", "0x"+hex.EncodeToString(data))
	q.Set("height", fmt.Sprint(c.Height))
	u := c.RPC + "/abci_query?" + q.Encode()
	var lastErr error
	for attempt := 0; attempt < 5; attempt++ {
		c.mu.Lock()
		c.Calls++
		if c.Progress != nil {
			fmt.Fprintf(c.Progress, "query #%d %s %x\n", c.Calls, path, data)
		}
		c.mu.Unlock()
		resp, err := c.HTTP.Get(u)
		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(attempt+1) * time.Second)
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			continue
		}
		var r abciResp
		if err := json.Unmarshal(body, &r); err != nil {
			lastErr = fmt.Errorf("decode %s: %w (%.200s)", path, err, body)
			time.Sleep(time.Duration(attempt+1) * time.Second)
			continue
		}
		if r.Error != nil {
			return nil, fmt.Errorf("%s: %s %s", path, r.Error.Message, r.Error.Data)
		}
		if r.Result.Response.Code != 0 {
			return nil, fmt.Errorf("%s: code %d %s", path, r.Result.Response.Code, r.Result.Response.Log)
		}
		return r.Result.Response.Value, nil
	}
	return nil, lastErr
}

// Key reads one key of a module store at the client height. Unlike "/subspace"
// (which the IAVL store serves from the LATEST version whatever the height), "/key"
// is versioned, so this is the only historical read to rely on.
func (c *Client) Key(store string, key []byte) ([]byte, error) {
	return c.query("/store/"+store+"/key", key)
}

// KeysBatch reads many keys at the client height in one JSON-RPC batch POST
// (versioned "/key" reads). Missing keys come back as nil.
func (c *Client) KeysBatch(stores []string, keys [][]byte) ([][]byte, error) {
	type params struct {
		Path   string `json:"path"`
		Data   string `json:"data"`
		Height string `json:"height"`
		Prove  bool   `json:"prove"`
	}
	type req struct {
		JSONRPC string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Method  string `json:"method"`
		Params  params `json:"params"`
	}
	reqs := make([]req, len(keys))
	for i := range keys {
		reqs[i] = req{"2.0", i, "abci_query", params{"/store/" + stores[i] + "/key", hex.EncodeToString(keys[i]), fmt.Sprint(c.Height), false}}
	}
	body, _ := json.Marshal(reqs)
	var lastErr error
	for attempt := 0; attempt < 5; attempt++ {
		c.mu.Lock()
		c.Calls += len(keys)
		if c.Progress != nil {
			fmt.Fprintf(c.Progress, "batch of %d, total %d\n", len(keys), c.Calls)
		}
		c.mu.Unlock()
		resp, err := c.HTTP.Post(strings.TrimSuffix(c.RPC, "/")+"/", "application/json", bytes.NewReader(body))
		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(attempt+1) * time.Second)
			continue
		}
		raw, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			continue
		}
		var out []struct {
			ID     int `json:"id"`
			Result struct {
				Response struct {
					Code   int    `json:"code"`
					Log    string `json:"log"`
					Value  []byte `json:"value"`
					Height string `json:"height"`
				} `json:"response"`
			} `json:"result"`
			Error *struct {
				Message string `json:"message"`
				Data    string `json:"data"`
			} `json:"error"`
		}
		if err := json.Unmarshal(raw, &out); err != nil {
			lastErr = fmt.Errorf("decode batch: %w (%.200s)", err, raw)
			time.Sleep(time.Duration(attempt+1) * time.Second)
			continue
		}
		if len(out) != len(keys) {
			lastErr = fmt.Errorf("batch returned %d results for %d keys", len(out), len(keys))
			continue
		}
		vals := make([][]byte, len(keys))
		for _, o := range out {
			if o.Error != nil {
				return nil, fmt.Errorf("batch item %d: %s %s", o.ID, o.Error.Message, o.Error.Data)
			}
			if o.Result.Response.Code != 0 {
				return nil, fmt.Errorf("batch item %d: code %d %s", o.ID, o.Result.Response.Code, o.Result.Response.Log)
			}
			if o.Result.Response.Height != fmt.Sprint(c.Height) {
				return nil, fmt.Errorf("batch item %d answered at height %s", o.ID, o.Result.Response.Height)
			}
			vals[o.ID] = o.Result.Response.Value
		}
		return vals, nil
	}
	return nil, lastErr
}

// Store is a read-through KVStore for one module store.
type Store struct {
	c      *Client
	name   string
	mu     sync.Mutex
	db     *dbm.MemDB      // every pair fetched so far
	absent map[string]bool // keys known not to exist
	loaded [][]byte        // prefixes fetched completely
	// Err records the first remote failure; the store returns empty results after it.
	Err error
}

func New(c *Client, name string) *Store {
	return &Store{c: c, name: name, db: dbm.NewMemDB(), absent: map[string]bool{}}
}

func (s *Store) fail(err error) {
	if s.Err == nil {
		s.Err = fmt.Errorf("store %s: %w", s.name, err)
	}
}

func (s *Store) covered(key []byte) bool {
	for _, p := range s.loaded {
		if bytes.HasPrefix(key, p) {
			return true
		}
	}
	return false
}

// LoadPrefix fetches every pair under prefix (an empty prefix loads the whole
// store, one leading byte at a time, since the node rejects empty subspace data).
func (s *Store) LoadPrefix(prefix []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.loadPrefix(prefix)
}

// LoadPrefixSplit fetches every pair under prefix as 256 sub-prefix queries run by
// `workers` goroutines; large subspaces time out as a single query.
func (s *Store) LoadPrefixSplit(prefix []byte, workers int) error {
	type res struct {
		pairs kv.Pairs
		err   error
	}
	jobs := make(chan int)
	results := make(chan res, 256)
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for b := range jobs {
				sub := append(append([]byte{}, prefix...), byte(b))
				bz, err := s.c.query("/store/"+s.name+"/subspace", sub)
				var r res
				if err != nil {
					r.err = err
				} else if err := r.pairs.Unmarshal(bz); err != nil {
					r.err = fmt.Errorf("unmarshal subspace %x: %w", sub, err)
				}
				results <- r
			}
		}()
	}
	go func() {
		for b := 0; b < 256; b++ {
			jobs <- b
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()
	s.mu.Lock()
	defer s.mu.Unlock()
	var firstErr error
	for r := range results {
		if r.err != nil {
			if firstErr == nil {
				firstErr = r.err
			}
			continue
		}
		for _, p := range r.pairs.Pairs {
			_ = s.db.Set(p.Key, p.Value)
		}
	}
	if firstErr != nil {
		return firstErr
	}
	s.loaded = append(s.loaded, append([]byte{}, prefix...))
	return nil
}

func (s *Store) loadPrefix(prefix []byte) error {
	if s.covered(prefix) && len(prefix) > 0 {
		return nil
	}
	if len(prefix) == 0 {
		for b := 0; b < 256; b++ {
			if err := s.loadPrefix([]byte{byte(b)}); err != nil {
				return err
			}
		}
		s.loaded = append(s.loaded, []byte{})
		return nil
	}
	bz, err := s.c.query("/store/"+s.name+"/subspace", prefix)
	if err != nil {
		return err
	}
	var pairs kv.Pairs
	if err := pairs.Unmarshal(bz); err != nil {
		return fmt.Errorf("unmarshal subspace %x: %w", prefix, err)
	}
	for _, p := range pairs.Pairs {
		if err := s.db.Set(p.Key, p.Value); err != nil {
			return err
		}
	}
	s.loaded = append(s.loaded, append([]byte{}, prefix...))
	return nil
}

func (s *Store) get(key []byte) []byte {
	if v, _ := s.db.Get(key); v != nil {
		return v
	}
	if s.covered(key) || s.absent[string(key)] {
		return nil
	}
	v, err := s.c.query("/store/"+s.name+"/key", key)
	if err != nil {
		s.fail(err)
		return nil
	}
	if len(v) == 0 {
		s.absent[string(key)] = true
		return nil
	}
	_ = s.db.Set(key, v)
	return v
}

func (s *Store) Get(key []byte) []byte {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.get(key)
}

func (s *Store) Has(key []byte) bool { return s.Get(key) != nil }

func (s *Store) Set(key, value []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = s.db.Set(key, value)
	delete(s.absent, string(key))
}

func (s *Store) Delete(key []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = s.db.Delete(key)
	s.absent[string(key)] = true
}

// iterPrefix is the longest common prefix of the iteration domain.
func iterPrefix(start, end []byte) []byte {
	if start == nil || end == nil {
		return nil
	}
	n := 0
	for n < len(start) && n < len(end) && start[n] == end[n] {
		n++
	}
	return start[:n]
}

func (s *Store) ensureDomain(start, end []byte) {
	if err := s.loadPrefix(iterPrefix(start, end)); err != nil {
		s.fail(err)
	}
}

func (s *Store) Iterator(start, end []byte) types.Iterator {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ensureDomain(start, end)
	return s.snapshot(start, end, false)
}

func (s *Store) ReverseIterator(start, end []byte) types.Iterator {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ensureDomain(start, end)
	return s.snapshot(start, end, true)
}

// snapshot copies the domain so callers may write while iterating.
func (s *Store) snapshot(start, end []byte, reverse bool) types.Iterator {
	it, err := s.db.Iterator(start, end)
	if err != nil {
		s.fail(err)
		return newSliceIter(nil, start, end)
	}
	var pairs []kv.Pair
	for ; it.Valid(); it.Next() {
		pairs = append(pairs, kv.Pair{Key: append([]byte{}, it.Key()...), Value: append([]byte{}, it.Value()...)})
	}
	it.Close()
	if reverse {
		sort.Slice(pairs, func(i, j int) bool { return bytes.Compare(pairs[i].Key, pairs[j].Key) > 0 })
	}
	return newSliceIter(pairs, start, end)
}

func (s *Store) GetStoreType() types.StoreType { return types.StoreTypeIAVL }

func (s *Store) CacheWrap() types.CacheWrap { return cachekv.NewStore(s) }

func (s *Store) CacheWrapWithTrace(w io.Writer, tc types.TraceContext) types.CacheWrap {
	return cachekv.NewStore(tracekv.NewStore(s, w, tc))
}

type sliceIter struct {
	pairs      []kv.Pair
	i          int
	start, end []byte
}

func newSliceIter(p []kv.Pair, start, end []byte) *sliceIter {
	return &sliceIter{pairs: p, start: start, end: end}
}

func (it *sliceIter) Domain() ([]byte, []byte) { return it.start, it.end }
func (it *sliceIter) Valid() bool              { return it.i < len(it.pairs) }
func (it *sliceIter) Next()                    { it.i++ }
func (it *sliceIter) Key() []byte              { return it.pairs[it.i].Key }
func (it *sliceIter) Value() []byte            { return it.pairs[it.i].Value }
func (it *sliceIter) Error() error             { return nil }
func (it *sliceIter) Close() error             { return nil }

var _ types.KVStore = (*Store)(nil)
