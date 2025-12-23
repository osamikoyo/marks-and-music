package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/osamikoyo/music-and-marks/services/api/server"
)

func main() {
	configPath := "api-config.yaml"
	for i, arg := range os.Args {
		if arg == "--config" {
			configPath = os.Args[i+1]
		}
	}

	server, err := server.SetupApiServer(configPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	go func() {
		<-ctx.Done()

		server.Close(ctx)
	}()

	if err := server.Start(ctx); err != nil {
		log.Fatal(err)
	}
}

