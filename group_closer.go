package app

import (
	"context"
	"fmt"
	"io"
	"log/slog"
)

type CloserGroup []io.Closer

func (cg CloserGroup) Close(ctx context.Context, l *slog.Logger) {
	if len(cg) == 0 {
		return
	}

	waitCh := make(chan struct{}, 1)

	go func() {
		defer close(waitCh)

		for _, c := range cg {
			select {
			case <-ctx.Done():
				return
			default:
			}

			l.Info(fmt.Sprintf("Closing [%T]", c))

			err := c.Close()
			if err != nil {
				l.Error(fmt.Sprintf("Error on close [%T]: %s", c, err))
			}

			l.Info(fmt.Sprintf("[%T] closed", c))
		}
	}()

	select {
	case <-ctx.Done():
		l.Error("Not all entities were closed")
	case <-waitCh:
	}
}
