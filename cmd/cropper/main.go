package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/seggga/observability/internal/app/endpoint"
	"github.com/seggga/observability/internal/app/service"
	storage "github.com/seggga/observability/internal/pkg/storage/memory"
)

func main() {
	log.Print("Starting the app")

	// port := flag.String("port", "8000", "Port")
	// storageName := flag.String("storage", "storage.json", "data storage")
	// shutdownTimeout := flag.Int64("shutdown_timeout", 3, "shutdown timeout")
	// flag.Parse()

	//	  repo := repomem.New()
	//    repo := repopg.New()
	repo := storage.New()
	svc := service.New(repo)

	app := http.Server{
		Addr:    net.JoinHostPort("localhost", "8080"),
		Handler: endpoint.RegisterPublicHTTP(svc),
	}
	go func() {
		if err := app.ListenAndServe(); err != nil {
			log.Fatalf("listen and serve err: %v", err)
		}
	}()

	interrupt := make(chan os.Signal, 3)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	log.Print("Started app")

	sig := <-interrupt

	log.Printf("Sig: %v, stopping app", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := app.Shutdown(ctx); err != nil {
		log.Printf("shutdown err: %v", err)
	}
}
