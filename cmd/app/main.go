package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GrishaSkurikhin/WB_Task_L0/internal/broker/nats"
	orderCache "github.com/GrishaSkurikhin/WB_Task_L0/internal/cache/map-cache"
	"github.com/GrishaSkurikhin/WB_Task_L0/internal/config"
	"github.com/GrishaSkurikhin/WB_Task_L0/internal/lib/logger/sl"
	"github.com/GrishaSkurikhin/WB_Task_L0/internal/lib/logger/slogpretty"
	"github.com/GrishaSkurikhin/WB_Task_L0/internal/orders"
	restserver "github.com/GrishaSkurikhin/WB_Task_L0/internal/rest-server"
	orderStorage "github.com/GrishaSkurikhin/WB_Task_L0/internal/storage/postgresql"
	"github.com/GrishaSkurikhin/WB_Task_L0/pkg/closer"
	"golang.org/x/exp/slog"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"

	natsClientID    = "1"
	shutdownTimeout = 5 * time.Second
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info(
		"starting app",
		slog.String("env", cfg.Env),
		slog.String("version", "1"),
	)
	log.Debug("debug messages are enabled")

	// connect to postgresql
	storage, err := orderStorage.New(cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.User, cfg.Storage.Password, cfg.Storage.DBName)
	if err != nil {
		log.Error("failed to connect to postgresql", sl.Err(err))
		return
	}
	log.Info("successful connect to postgresql")

	// create cache
	cache := orderCache.New()
	log.Info("successful create cache")

	// add orders to cache
	err = orders.AddAllToCache(cache, storage)
	if err != nil {
		log.Error("failed to add orders to cache", sl.Err(err))
		return
	}
	log.Info("all orders have been added to the cache")

	// create nats client
	natsClient, err := nats.New(cfg.Nats.Host, cfg.Nats.Port, cfg.Nats.Cluster, natsClientID)
	if err != nil {
		log.Error("failed to create nats client", sl.Err(err))
		return
	}

	// subscribe on channel
	err = natsClient.SubscribeOrderChannel(log, cache, storage, cfg.Nats.SubjectOrderAdd)
	if err != nil {
		log.Error("failed to subscribe nats channel", sl.Err(err))
		return
	}

	// create server
	srv, err := restserver.New(cfg, log, cache)
	if err != nil {
		log.Error("failed to create server", sl.Err(err))
		return
	}
	log.Info("starting server", slog.String("address", cfg.RestServer.Address))

	c := &closer.Closer{}
	c.Add(srv.Close)
	c.Add(natsClient.Disconnect)
	c.Add(storage.Disconnect)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.Start(); err != nil {
			log.Error("failed to start server", sl.Err(err))
		}
	}()

	log.Info("server started")

	<-ctx.Done()
	log.Info("stopping server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := c.Close(shutdownCtx); err != nil {
		log.Error("closer error", sl.Err(err))
	}

	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
