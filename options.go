package app

import (
	"io"
	"log/slog"
)

type Options struct {
	servers          []Server
	gracefulServices []GracefulService
	closers          []io.Closer

	logger *slog.Logger
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	var o = Options{
		logger: slog.Default(),
	}

	for _, opt := range opts {
		opt(&o)
	}

	return o
}

func WithServer(srvr ...Server) Option {
	return func(o *Options) {
		o.servers = append(o.servers, srvr...)
	}
}
