package transport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalAndUnmarshal(t *testing.T) {
	data := []byte("some infomoation here")
	input_t := TCPCommand{
		Command: byte(0x2),
		Data:    data,
	}
    marshalled_data, err := input_t.MarshalBinary()
    if err != nil {
        t.Error(err)
    }
    
    output_t := TCPCommand{}
    
    output_t.UnmarshalBinary(marshalled_data)

    assert.Equal(t, byte(0x2), output_t.Command)
    assert.Equal(t, data, output_t.Data)
}
