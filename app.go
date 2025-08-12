package app

import (
	"context"
	"errors"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/nuvotlyuba/app/signal"
)

var ErrStopped = errors.New("application stopped")

const shutdownTimeout = 10 * time.Second

type Application struct {
	servers          ServerGroup
	gracefulServices GracefulServiceGroup
	closers          CloserGroup

	stopped atomic.Bool

	logger *slog.Logger
}

func New(cfg Config, opt ...Option) *Application {
	o := newOptions(opt...)

	var app = &Application{
		servers:          o.servers,
		gracefulServices: o.gracefulServices,
		closers:          o.closers,
		logger:           o.logger,
	}

	return app
}

func (a *Application) Run() {
	if a.stopped.Load() {
		panic(ErrStopped)
	}

	a.logger.Info("Starting application...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a.gracefulServices.Run(ctx, a.logger)
	a.servers.Serve(a.logger)

	a.logger.Info("Application started")

	<-signal.Watch()
	a.stop(ctx)
}

func (a *Application) stop(ctx context.Context) {
	defer a.stopped.Store(true)

	a.logger.Info("Stopping application...")

	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	a.servers.Shutdown(ctx, a.logger)
	a.gracefulServices.Stop(ctx, a.logger)
	a.closers.Close(ctx, a.logger)

	a.logger.Info("Application stopped")

}
