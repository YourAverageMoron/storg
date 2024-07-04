package raft

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeoutWaits(t *testing.T) {
	d := func() time.Duration {
		return time.Duration(1) * time.Second
	}
	timeout := NewTimeout(d)
	timeout.Start()

	select {
	case <-timeout.Consume():
		assert.True(t, true)
	case <-time.After(time.Second * time.Duration(2)):
		assert.True(t, false)
	}

	select {
	case <-timeout.Consume():
		assert.True(t, false)
	case <-time.After(time.Millisecond * time.Duration(10)):
		assert.True(t, true)
	}

    timeout.Stop()

}


