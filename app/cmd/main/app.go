package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/GermanBogatov/TodoApp/app/internal/config"
	"github.com/GermanBogatov/TodoApp/app/internal/handler"
	"github.com/GermanBogatov/TodoApp/app/internal/service"
	"github.com/GermanBogatov/TodoApp/app/internal/storage"
	"github.com/GermanBogatov/TodoApp/app/pkg/jwt"
	"github.com/GermanBogatov/TodoApp/app/pkg/logging"
	"github.com/GermanBogatov/TodoApp/app/pkg/postgresql"
	"github.com/GermanBogatov/TodoApp/app/pkg/redis"
	"github.com/GermanBogatov/TodoApp/app/pkg/shutdown"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"syscall"
	"time"
)

func main() {

	logging.Init()
	logger := logging.GetLogger()
	logger.Println("logger initialized")

	logger.Println("Postgresql config initializing")
	cfg := config.GetConfig()

	logger.Println("Redis initializing")
	RedisClient, err := redis.NewClient(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)

	logger.Println("JWT Helper initializing")
	NewHelper := jwt.NewHelper(logger, RedisClient)

	logger.Println("Postgresql client initializing")
	PostgresqlClient, err := postgresql.NewClient(context.Background(), 5, cfg.PostgresqlDB.Username, cfg.PostgresqlDB.Password,
		cfg.PostgresqlDB.Host, cfg.PostgresqlDB.Port, cfg.PostgresqlDB.Database)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Println("Storage initializing")
	Storage := storage.NewRepository(PostgresqlClient, logger)
	if err != nil {
		panic(err)
	}

	logger.Println("Service initializing")
	Service, err := service.NewService(Storage, logger)
	if err != nil {
		panic(err)
	}

	logger.Println("Handler initializing")

	Handler, err := handler.NewHandler(Service, logger, NewHelper)
	if err != nil {
		panic(err)
	}

	logger.Println("start application")
	start(Handler.InitRoutes(), logger, cfg)

}

func start(router http.Handler, logger logging.Logger, cfg *config.Config) {
	var server *http.Server
	var listener net.Listener

	if cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		socketPath := path.Join(appDir, "app.sock")
		logger.Infof("socket path: %s", socketPath)

		logger.Info("create and listen unix socket")
		listener, err = net.Listen("unix", socketPath)
		if err != nil {
			logger.Fatal(err)
		}
	} else {
		logger.Infof("bind application to host: %s and port: %s", cfg.Listen.BindIP, cfg.Listen.Port)

		var err error

		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		if err != nil {
			logger.Fatal(err)
		}
	}

	server = &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go shutdown.Graceful([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM}, server)

	logger.Println("application initialized and started")

	if err := server.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logger.Warn("server shutdown")
		default:
			logger.Fatal(err)
		}
	}
}
