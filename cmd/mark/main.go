package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/osamikoyo/music-and-marks/services/mark/app"
)

func main() {
	configpath := ""

	for i, arg := range os.Args {
		if arg == "--config" {
			configpath = os.Args[i+1]
		}
	}

	app, err := app.SetupApp(configpath)
	if err != nil {
		log.Fatal(err)

		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := app.Run(ctx); err != nil {
		log.Fatal(err)

		return
	}
}
