package signals

import (
	"context"
	"log/slog"
	"time"

	"github.com/lgtmpub/smartchat/pkg/logger"
)

var defaultShutdownTimeout = 3 * time.Second

type Shutdown struct {
	serverShutdownTimeout time.Duration
	logger                logger.Logger
}

type Option func(*Shutdown)

func WithServerShutdownTimeout(timeout time.Duration) Option {
	return func(s *Shutdown) {
		s.serverShutdownTimeout = timeout
	}
}

func WithLogger(logger logger.Logger) Option {
	return func(s *Shutdown) {
		s.logger = logger
	}
}

func NewShutdown(opts ...Option) (*Shutdown, error) {
	srv := &Shutdown{
		serverShutdownTimeout: defaultShutdownTimeout,
	}
	for _, opt := range opts {
		opt(srv)
	}

	return srv, nil
}

type Service interface {
	PreShutdown()
	Shutdown(context.Context) error
}

func (s *Shutdown) Graceful(stopCh <-chan struct{}, srvs ...Service) {
	// wait for SIGTERM or SIGINT
	<-stopCh
	ctx, cancel := context.WithTimeout(context.Background(), s.serverShutdownTimeout)
	defer cancel()
	s.logger.Info("shutting down server", slog.String("timeout", s.serverShutdownTimeout.String()))

	for _, srv := range srvs {
		srv.PreShutdown()
	}

	// There could be a period where a terminating pod may still receive requests. Implementing a brief wait can mitigate this.
	// See: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-termination
	// the readiness check interval must be lower than the timeout
	time.Sleep(s.serverShutdownTimeout)

	for _, srv := range srvs {
		if err := srv.Shutdown(ctx); err != nil {
			s.logger.Error("shutdown function failed", slog.String("error", err.Error()))
		}
	}
}
