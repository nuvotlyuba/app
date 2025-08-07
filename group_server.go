package app

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

type Server interface {
	Serve() error
	Shutdown(context.Context) error
}

type ServerGroup []Server

func (sg ServerGroup) Serve(l *slog.Logger) {
	l.Info("Starting servers...")

	var wg sync.WaitGroup
	wg.Add(len(sg))

	for _, server := range sg {
		go func(s Server) {
			defer wg.Done()

			err := s.Serve()
			if err != nil {
				panic(err)
			}
		}(server)
	}

	wg.Wait()
}

func (sg ServerGroup) Shutdown(ctx context.Context, l *slog.Logger) {
	if len(sg) == 0 {
		return
	}

	l.Info("Stopping servers...")

	var (
		wg     sync.WaitGroup
		waitCh = make(chan struct{}, 1)
	)
	wg.Add(len(sg))

	go func() {
		wg.Wait()
		close(waitCh)
	}()

	for _, server := range sg {
		go func(s Server) {
			defer wg.Done()

			err := s.Shutdown(ctx)
			if err != nil {
				l.Error(fmt.Sprintf("Error on stopping server [%T]: %s", s, err))
			}
		}(server)
	}

	select {
	case <-ctx.Done():
		l.Error("Not all servers were finished")
	case <-waitCh:
	}

}
