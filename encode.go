package plist

import "io"

// Encoder ...
type Encoder struct {
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{}
}

// Encode ...
func (enc *Encoder) Encode(v interface{}) error {
	return nil
}
