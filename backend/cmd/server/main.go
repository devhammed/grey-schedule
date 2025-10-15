package main

import (
	"log"
	"os"

	"github.com/devhammed/grey-schedule/backend/internal/grpcapi"
	"github.com/devhammed/grey-schedule/backend/internal/httpapi"
	"github.com/devhammed/grey-schedule/backend/internal/store"
)

func main() {
	httpAddr := getEnv("HTTP_ADDR", ":8080")
	grpcAddr := getEnv("GRPC_ADDR", ":8081")
	dsn := getEnv("POSTGRES_URL", "postgres://root:@localhost:5432/grey_schedule?sslmode=disable")
	st, err := store.NewPostgresStore(dsn)

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		api := grpcapi.NewServer(st)

		log.Printf("gRPC server listening on %s", grpcAddr)

		if err := api.Start(grpcAddr); err != nil {
			log.Printf("gRPC server stopped: %v", err)
		}
	}()

	go func() {
		api := httpapi.NewServer(st)

		log.Printf("HTTP server listening on %s", httpAddr)

		if err := api.Start(httpAddr); err != nil {
			log.Printf("HTTP server stopped: %v", err)
		}
	}()

	select {}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return def
}
