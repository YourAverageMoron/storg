package transport

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestMarshalAndUnmarshal(t *testing.T) {
	data := []byte("some infomoation here")
    message := Message{
        Command: RegisterPeer,
		Payload: data,
	}
    input_t := TCPMessage{Message: message}
	marshalled_data, err := input_t.MarshalBinary()
	if err != nil {
		t.Error(err)
	}

	output_t := TCPMessage{}
    
    r := bytes.NewReader(marshalled_data)
	output_t.UnmarshalBinary(r)
	assert.Equal(t, RegisterPeer, output_t.Command)
	assert.Equal(t, data, output_t.Payload)
}
