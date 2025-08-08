package signal

import (
	"os"
	"os/signal"
	"syscall"
)

func Watch() <-chan os.Signal {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	return interrupt
}
