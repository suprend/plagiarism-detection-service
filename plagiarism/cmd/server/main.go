package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"plagiarism/internal/api/http/router"
	"plagiarism/internal/application/usecase"
	"plagiarism/internal/domain"
	"plagiarism/internal/infrastructure/config"
	"plagiarism/internal/infrastructure/filestorage"
	"plagiarism/internal/infrastructure/report"
	"plagiarism/internal/infrastructure/worker"
)

func main() {
	reportStore := report.NewFileReportStore("plagiarism/reports")
	fsClient := filestorage.NewClient(config.FilestorageURL())
	w := worker.NewWorker(reportStore, fsClient, config.MatchThreshold(), config.WorkerCount(), func(rep domain.CheckReport, err error) {
		log.Printf("failed to save report work=%s submission=%s: %v", rep.WorkID, rep.SubmissionID, err)
	})
	checkUseCase := usecase.NewCheckService(reportStore, w)

	r := router.NewRouter(checkUseCase)
	handler := r.SetupRoutes()

	port := config.ServerPort()
	addr := ":" + port
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		log.Printf("plagiarism service starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down plagiarism service...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	w.Close()
	log.Println("plagiarism service stopped")
}
