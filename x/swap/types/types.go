package types

import "encoding/hex"

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

type Secret []byte

func (s Secret) MarshalJSON() ([]byte, error) {
	return []byte("\"" + hex.EncodeToString(s) + "\""), nil
}

func (s *Secret) UnmarshalJSON(b []byte) error {
	decoded, err := hex.DecodeString(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}
	*s = decoded
	return nil
}

func (s Secret) Size() int {
	raw, _ := s.MarshalJSON()

	return len(raw)
}

func (s Secret) MarshalTo(bytes []byte) ([]byte, error) {
	// todo
	return bytes, nil
}

func (s Secret) Unmarshal(bytes []byte) error {
	return s.UnmarshalJSON(bytes)
}
