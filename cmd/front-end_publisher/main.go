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

func init() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "../../pub-sub-359008-ff94c59da4aa.json")
}

func main() {
	cfg := config.Set()

	r := gin.Default()
	handler := handler.New(cfg)
	router := handler.Routes(r)

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
