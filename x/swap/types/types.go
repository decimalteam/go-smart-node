package types

import "encoding/hex"

////////////////////////////////////////////////////////////////
// Hash
////////////////////////////////////////////////////////////////

type Hash [32]byte

func (h Hash) MarshalJSON() ([]byte, error) {
	return []byte("\"" + hex.EncodeToString(h[:]) + "\""), nil
}

func (h *Hash) UnmarshalJSON(b []byte) error {
	decoded, err := hex.DecodeString(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}
	copy(h[:], decoded)
	return nil
}

func (h *Hash) Size() int {
	s, _ := h.MarshalJSON()
	return len(s)
}

func (h Hash) MarshalTo(bytes []byte) ([]byte, error) {
	bytes, err := h.MarshalJSON()

	return bytes, err
}

func (h Hash) Unmarshal(bytes []byte) error {
	return h.UnmarshalJSON(bytes)
}

func (h *Hash) Copy() *Hash {
	var result Hash
	copy(result[:], (*h)[:])
	return &result
}

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

////////////////////////////////////////////////////////////////
// Chain
////////////////////////////////////////////////////////////////

func NewChain(number uint32, name string, active bool) Chain {
	return Chain{Number: number, Name: name, Active: active}
}
