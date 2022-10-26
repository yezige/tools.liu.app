//go:build !lambda
// +build !lambda

package router

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-queue/queue"
	"github.com/yezige/tools.liu.app/config"
	"github.com/yezige/tools.liu.app/logx"
	"golang.org/x/sync/errgroup"
)

// RunHTTPServer provide run http or https protocol.
func RunHTTPServer(ctx context.Context, cfg *config.ConfYaml, q *queue.Queue, s ...*http.Server) (err error) {
	var server *http.Server

	if !cfg.Core.Enabled {
		logx.LogAccess.Info("httpd server is disabled.")
		return nil
	}

	if len(s) == 0 {
		server = &http.Server{
			Addr:    cfg.Core.Address + ":" + cfg.Core.Port,
			Handler: routerEngine(cfg, q),
		}
	} else {
		server = s[0]
	}

	logx.LogAccess.Info("HTTPD server is running on " + cfg.Core.Port + " port.")

	return startServer(ctx, server, cfg)
}

func listenAndServe(ctx context.Context, s *http.Server, cfg *config.ConfYaml) error {
	var g errgroup.Group
	g.Go(func() error {
		<-ctx.Done()
		timeout := time.Duration(cfg.Core.ShutdownTimeout) * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		return s.Shutdown(ctx)
	})
	g.Go(func() error {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	return g.Wait()
}

func startServer(ctx context.Context, s *http.Server, cfg *config.ConfYaml) error {
	return listenAndServe(ctx, s, cfg)
}
