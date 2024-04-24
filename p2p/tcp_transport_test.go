package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listenAddr := ":4000"
	tr := NewTCPTransport(listenAddr).(*TCPTransport)
	assert.Equal(t, tr.listenAddress, listenAddr)

	// Server
	assert.Nil(t, tr.ListenAndAccept())
}
