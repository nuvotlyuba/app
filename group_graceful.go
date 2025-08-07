package app

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type GracefulService interface {
	Run(context.Context) error
	Stop(context.Context) error
}

type GracefulServiceGroup []GracefulService

func (gsg GracefulServiceGroup) Run(ctx context.Context, l *slog.Logger) {
	var wg sync.WaitGroup
	wg.Add(len(gsg))

	for _, service := range gsg {
		go func(s GracefulService) {
			defer wg.Done()
			runGracefulService(ctx, l, s)
		}(service)
	}

	wg.Wait()
}

func runGracefulService(ctx context.Context, l *slog.Logger, s GracefulService) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		l.Info(fmt.Sprintf("Starting graceful service [%T]", s))

		err := s.Run(ctx)
		if err == nil {
			return
		}

		l.Error(fmt.Sprintf("Graceful service [%T] finished with error: %s"), s, err)
		time.Sleep(time.Second)
	}
}

func (sgs GracefulServiceGroup) Stop(ctx context.Context, l *slog.Logger) {
	if len(sgs) == 0 {
		return
	}

	waitCh := make(chan struct{}, 1)

	go func() {
		defer close(waitCh)

		for _, s := range sgs {
			select {
			case <-ctx.Done():
				return
			default:
			}

			l.Info(fmt.Sprintf("Stopping graceful service [%T]", s))
			err := s.Stop(ctx)
			if err != nil {
				l.Error(fmt.Sprintf("Error stopping [%T]: %s", s, err))
			}
			l.Info(fmt.Sprintf("Graceful service [%T] stopped", s))
		}
	}()

	select {
	case <-ctx.Done():
		l.Error("Not all graceful services were finished")
	case <-waitCh:
	}
}
