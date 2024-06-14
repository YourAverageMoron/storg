package transport

import (
	"encoding/gob"
	"io"
)

type Encoder interface {
	Encode(w io.Writer, payload any) error
	Decode(r io.Reader, payload any) error
    Register(...any)
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

func(_ GobEncoder) Register(payloads ...any) {
    for _, p := range payloads {
	    gob.Register(p)
    }

}
