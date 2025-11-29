package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"filestorage/internal/api/http/router"
	"filestorage/internal/application/usecase"
	"filestorage/internal/infrastructure/config"
	"filestorage/internal/infrastructure/repository/postgres"
	"filestorage/internal/infrastructure/repository/s3"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	dbConfig := config.LoadDatabaseConfig()
	s3Config := config.LoadS3Config()

	pool, err := pgxpool.New(ctx, dbConfig.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	submissionRepo := postgres.NewPostgresRepository(pool)
	s3Repo, err := s3.NewS3Repository(ctx, s3Config.Bucket, s3Config.Endpoint, s3Config.Region)
	if err != nil {
		log.Fatalf("Failed to initialize S3 repository: %v", err)
	}

	submitUseCase := usecase.NewSubmitUseCase(submissionRepo, s3Repo)
	getSubmissionsUseCase := usecase.NewGetSubmissionsUseCase(submissionRepo)
	downloadSubmissionUseCase := usecase.NewDownloadSubmissionUseCase(submissionRepo, s3Repo)

	r := router.NewRouter(submitUseCase, getSubmissionsUseCase, downloadSubmissionUseCase)
	handler := r.SetupRoutes()

	port := ":" + config.ServerPort()
	srv := &http.Server{
		Addr:    port,
		Handler: handler,
	}

	go func() {
		fmt.Printf("Server starting on port %s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server exited")
}
