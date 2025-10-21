package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/devhammed/grey-schedule/backend/internal/grpcapi"
	"github.com/devhammed/grey-schedule/backend/internal/httpapi"
	"github.com/devhammed/grey-schedule/backend/internal/store"
)

func main() {
	grpcAddr := getEnv("GRPC_ADDR", ":8081")

	httpAddr := getEnv("HTTP_ADDR", ":8080")

	storeType := getEnv("STORE_TYPE", "postgres")

	var (
		st  store.Store
		err error
	)

	if storeType == "postgres" {
		dsn := getEnv("POSTGRES_URL", "postgres://root:@localhost:5432/grey_schedule?sslmode=disable")

		st, err = store.NewPostgresStore(dsn)

		if err != nil {
			log.Fatalf("failed to init store: %v", err)
		}
	} else if storeType == "in-memory" {
		st = store.NewInMemoryStore()
	} else {
		log.Fatalf("unknown store type: %s", storeType)
	}

	grpcServer := grpcapi.NewServer(st, grpcAddr)

	httpServer := httpapi.NewServer(st, httpAddr)

	go func() {
		log.Printf("gRPC server listening on %s", grpcAddr)

		if err := grpcServer.Start(); err != nil {
			log.Printf("gRPC server stopped: %v", err)
		}
	}()

	go func() {
		log.Printf("HTTP server listening on %s", httpAddr)

		if err := httpServer.Start(); err != nil {
			log.Printf("HTTP server stopped: %v", err)
		}
	}()

	processId := os.Getpid()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Servers running on PID %d", processId)

	<-quit

	log.Println("Shutting down gracefully, press Ctrl+C again to force")

	var wg sync.WaitGroup

	wg.Go(func() {
		if err := httpServer.Stop(); err != nil {
			log.Fatalf("HTTP Server forced to shutdown: %v", err)
		}

		log.Println("HTTP Server stopped")
	})

	wg.Go(func() {
		if err := grpcServer.Stop(); err != nil {
			log.Fatalf("GRPC Server forced to shutdown: %v", err)
		}

		log.Println("GRPC Server stopped")
	})

	wg.Wait()
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return def
}
