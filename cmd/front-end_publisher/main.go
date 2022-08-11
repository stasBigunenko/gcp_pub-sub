package main

import (
	"context"
	"gcp_pub-sub/cmd/front-end_publisher/handler"
	"log"
	"net/http"
	"os"
	"os/signal"

	"gcp_pub-sub/cmd/front-end_publisher/config"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Set()

	r := gin.Default()
	h := handler.New(cfg)
	router := h.Routes(r)

	srv := http.Server{
		Addr:    cfg.Port,
		Handler: router,
	}

	c := make(chan os.Signal, 1)
	defer close(c)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		srv.Shutdown(context.Background())
	}()

	log.Printf("HTTP server started on port: %v\n", cfg.Port)

	if err := srv.ListenAndServe(); err != nil {
		log.Printf("failed to serve:+%v\n", err)
	}
}
