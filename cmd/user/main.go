package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/osamikoyo/music-and-marks/services/user/app"
)

func main() {
	config_path := "user-service.yaml"

	for i, arg := range os.Args{
		if arg == "--config" {
			config_path = os.Args[i+1]
		}
	}

	app, err := app.SetupApp(config_path)
	if err != nil{
		log.Fatal(err)

		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	if err = app.Start(ctx);err != nil{
		log.Fatal(err)

		return
	}
}