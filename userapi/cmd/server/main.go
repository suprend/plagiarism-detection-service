package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"userapi/internal/api/http/router"
	"userapi/internal/application/usecase"
	"userapi/internal/infrastructure/config"
	"userapi/internal/infrastructure/filestorage"
	"userapi/internal/infrastructure/plagiarism"
	"userapi/internal/infrastructure/wordcloud"
)

func main() {
	fsClient := filestorage.NewClient(config.FilestorageURL())
	plagClient := plagiarism.NewService(plagiarism.NewClient(config.PlagiarismURL()))

	submitUseCase := usecase.NewSubmitUseCase(fsClient, plagClient)
	reportsUseCase := usecase.NewReportsUseCase(plagClient)
	wcClient := wordcloud.NewClient(config.WordcloudURL())
	wordcloudUseCase := usecase.NewWordcloudUseCase(fsClient, wcClient, config.WordcloudDir())

	r := router.NewRouter(submitUseCase, reportsUseCase, wordcloudUseCase)
	handler := r.SetupRoutes()

	port := ":" + config.ServerPort()
	srv := &http.Server{
		Addr:    port,
		Handler: handler,
	}

	go func() {
		log.Printf("userapi gateway starting on %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down userapi gateway...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("userapi gateway stopped")
}
