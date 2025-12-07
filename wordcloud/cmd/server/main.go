package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wordcloud/internal/api/http/router"
	"wordcloud/internal/application/usecase"
	"wordcloud/internal/infrastructure/config"
	"wordcloud/internal/infrastructure/filestorage"
	"wordcloud/internal/infrastructure/wordcloud"
)

func main() {
	fsClient := filestorage.NewClient(config.FilestorageURL())
	wcClient := wordcloud.NewClient(config.WordcloudGeneratorURL())
	wcService := usecase.NewWordcloudService(fsClient, wcClient, config.WordcloudDir())

	r := router.NewRouter(wcService)
	handler := r.SetupRoutes()

	port := config.ServerPort()
	addr := ":" + port
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		log.Printf("wordcloud service starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down wordcloud service...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("wordcloud service stopped")
}
