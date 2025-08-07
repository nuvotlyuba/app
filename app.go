package app

import (
	"log/slog"
	"sync/atomic"
)

type Application struct {
	servers          ServerGroup
	gracefulServices GracefulServiceGroup
	closers          CloserGroup

	stopped atomic.Bool

	logger *slog.Logger
}

func New(cfg Config, opt ...Options) *Application {
	o := newOptions(opt...)

	var app = &Application{}

	return app
}
