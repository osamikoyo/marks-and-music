package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/osamikoyo/music-and-marks/services/user/app"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	configpath := ""

	for i, arg := range os.Args {
		if arg == "--config" {
			configpath = os.Args[i+1]
		}
	}

	app, err := app.SetupApp(configpath)
	if err != nil {
		return
	}

	if err = app.Start(ctx); err != nil {
		return
	}
}
