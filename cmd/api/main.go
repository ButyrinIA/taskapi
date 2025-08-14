package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ButyrinIA/taskapi/internal/adapters/handlers"
	"github.com/ButyrinIA/taskapi/internal/adapters/repositories"
	"github.com/ButyrinIA/taskapi/internal/core/usecases"
)

func main() {
	//инит репо
	repo := repositories.NewInMemoryTaskRepository()

	//инит логер
	logChan := make(chan string, 100)
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Println("Logger goroutine shutting down")
				return
			case msg := <-logChan:
				log.Println(msg)
			}
		}
	}()

	//инит юзкейс
	taskUsecase := usecases.NewTaskUsecase(repo, logChan)

	//инит хендлер
	handler := handlers.NewTaskHandler(taskUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", handler.HandleTasks)
	mux.HandleFunc("/tasks/", handler.HandleTaskByID)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	//Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	wg.Wait()

	log.Println("Server exiting")
}
