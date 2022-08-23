package main

import (
	"Intern/gcp_pub-sub/modules/subscriber/pkg/app"
	"context"
	"log"
	"os"
	"os/signal"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx := context.Background()

	// read config
	application, err := app.Create()
	if err != nil {
		log.Fatalf("app create internal problem: %w\n", err)
		os.Exit(1)
	}

	// create shutdown
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	if err = application.Run(ctx); err != nil {
		log.Fatal("Problems with server run: ", err)
	}

}
