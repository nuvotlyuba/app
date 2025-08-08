package signal

import (
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnit_Watch(t *testing.T) {
	sigChan := Watch()
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

	assert.Eventually(t, func() bool {
		<-sigChan
		return true
	}, time.Second, 100*time.Millisecond)
}
