package transport

import (
	"encoding/gob"
	"io"
)

type Encoder interface {
	Encode(w io.Writer, payload any) error
	Decode(r io.Reader, payload any) error
}

type GobEncoder struct{}

func (_ GobEncoder) Encode(w io.Writer, payload any) error {
	enc := gob.NewEncoder(w)
	return enc.Encode(payload)
}

func (_ GobEncoder) Decode(r io.Reader, payload any) error {
	dec := gob.NewDecoder(r)
	return dec.Decode(payload)
}
