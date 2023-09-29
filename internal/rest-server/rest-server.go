package restserver

import (
	"context"
	"fmt"
	"net/http"

	orderCache "github.com/GrishaSkurikhin/WB_Task_L0/internal/cache/map-cache"
	"github.com/GrishaSkurikhin/WB_Task_L0/internal/config"
	"github.com/GrishaSkurikhin/WB_Task_L0/internal/orders"
	getorder "github.com/GrishaSkurikhin/WB_Task_L0/internal/rest-server/handlers/get-order"
	mwLogger "github.com/GrishaSkurikhin/WB_Task_L0/internal/rest-server/middleware/logger"
	orderStorage "github.com/GrishaSkurikhin/WB_Task_L0/internal/storage/postgresql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
)

type restServer struct {
	httpServer *http.Server
	storage    *orderStorage.OrderStorage
	cache      *orderCache.MapCache
}

func New(cfg *config.Config, log *slog.Logger) (*restServer, error) {
	const op = "restserver.New"

	storage, err := orderStorage.New(cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.User, cfg.Storage.Password, cfg.Storage.DBName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	cache := orderCache.New()

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/users", func(r chi.Router) {
		r.Get("/get", getorder.New(log, cache))
	})

	srv := &http.Server{
		Addr:         cfg.RestServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.RestServer.Timeout,
		WriteTimeout: cfg.RestServer.Timeout,
		IdleTimeout:  cfg.RestServer.IdleTimeout,
	}

	return &restServer{
		httpServer: srv,
		storage: storage,
		cache: cache,
	}, nil
}

func (srv *restServer) Start() error {
	const op = "restserver.Start"
	err := orders.AddAllToCache(srv.cache, srv.storage)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = srv.httpServer.ListenAndServe()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (srv *restServer) Close(ctx *context.Context) error {
	const op = "restserver.Close"
	srv.storage.Close()

	err := srv.httpServer.Shutdown(*ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
