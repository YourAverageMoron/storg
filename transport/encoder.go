package transport

import (
	"encoding/gob"
	"io"
)

type Payload struct{}

type Encoder interface {
	Encode(w io.Writer, payload any) (io.Writer, error)
	Decode(r io.Reader, payload any) error
}

type GobEncoder struct{}

func (_ GobEncoder) Encode(w io.Writer, payload any) (io.Writer, error) {
	enc := gob.NewEncoder(w)
	if err := enc.Encode(payload); err != nil {
		return nil, err
	}
	// TODO: THIS MIGHT NOT NEED TO RETURN THE WRITER
	return w, nil
}

func (_ GobEncoder) Decode(r io.Reader, payload any) error {
	dec := gob.NewDecoder(r)
	return dec.Decode(payload)
}
