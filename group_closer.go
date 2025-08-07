package app

import (
	"context"
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

	}
}
