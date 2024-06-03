package transport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestMarshalAndUnmarshal(t *testing.T) {
	data := []byte("some infomoation here")
    message := Message{
        Command: RegisterPeer,
		Data: data,
	}
    input_t := TCPMessage{Message: message}
	marshalled_data, err := input_t.MarshalBinary()
	if err != nil {
		t.Error(err)
	}

	output_t := TCPMessage{}

	output_t.UnmarshalBinary(marshalled_data)

	assert.Equal(t, RegisterPeer, output_t.Command)
	assert.Equal(t, data, output_t.Data)
}
