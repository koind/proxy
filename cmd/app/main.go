package main

import (
	"context"
	"errors"
	"github.com/koind/proxy/internal/config"
	"github.com/koind/proxy/internal/handler"
	"github.com/koind/proxy/internal/storage/memory"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/pflag"
)

var configPath string

const DefaultConfigPath = "config/config.toml"

func init() {
	pflag.StringVarP(&configPath, "config", "c", DefaultConfigPath, "Путь до конфигурационного файла")
}

// @title Proxy
// @version 1.0
// @description Simple service on Go for proxying HTTP requests to third-party services.

// @host localhost:8080
// @BasePath /

func main() {
	pflag.Parse()

	cfg, err := config.Init(configPath)
	if err != nil {
		log.Fatal(err)
	}

	repository := memory.NewRequestRepository()
	srv := handler.NewHTTPServer(repository, cfg.HTTPServer.GetDomain())

	go func() {
		log.Printf("Запуск сервера, %s", cfg.HTTPServer.GetDomain())

		if err := srv.Start(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("Ошибка при старте сервера: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Выключение сервера ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Сервер принудительно остановлен:", err)
	}

	log.Println("Сервер остановлен")
}
