package system

import "log/slog"

func New(cfg rest.Config, l *slog.Logger) *rest.RestServer {
	if cfg.Port <= 0 || cfg.Port == rest.DefaultPort {
		cfg.Port = 8090
	}

	return rest.New(cfg, []rest.Handler{metrics{}},
		rest.WithName("system"),
		rest.WithLogger(l))
}
