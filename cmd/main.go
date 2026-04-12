package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	app, err := wireApp(ctx)
	if err != nil {
		log.Fatal(err)
	}

	go app.run()
	app.watchForShutdown(ctx)
}
