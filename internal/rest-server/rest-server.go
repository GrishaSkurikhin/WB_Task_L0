package restserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/GrishaSkurikhin/WB_Task_L0/internal/config"
	"github.com/GrishaSkurikhin/WB_Task_L0/internal/orders"
	getorder "github.com/GrishaSkurikhin/WB_Task_L0/internal/rest-server/handlers/get-order"
	mwLogger "github.com/GrishaSkurikhin/WB_Task_L0/internal/rest-server/middleware/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
)

type restServer struct {
	*http.Server
}

func New(cfg *config.Config, log *slog.Logger, cache orders.CacheGetter) (*restServer, error) {
	const op = "restserver.New"

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/order", func(r chi.Router) {
		r.Get("/get", getorder.New(log, cache))
	})

	srv := &http.Server{
		Addr:         cfg.RestServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.RestServer.Timeout,
		WriteTimeout: cfg.RestServer.Timeout,
		IdleTimeout:  cfg.RestServer.IdleTimeout,
	}

	return &restServer{srv}, nil
}

func (srv *restServer) Start() error {
	const op = "restserver.Start"

	err := srv.ListenAndServe()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (srv *restServer) Close(ctx context.Context) error {
	const op = "restserver.Close"

	err := srv.Shutdown(ctx)
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
