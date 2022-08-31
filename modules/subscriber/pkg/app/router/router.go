package router

import (
	"Intern/gcp_pub-sub/modules/subscriber/handler"
	"Intern/gcp_pub-sub/modules/subscriber/pkg/app/config"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

// Router is the object that represent HTTP API layer
type Router struct {
	serverPort *config.HTTPServerConfiguration
	handler    handler.Handler
}

// New is a constructor of the Router entity
func New(httpPort *config.HTTPServerConfiguration, event handler.Handler) *Router {
	return &Router{
		serverPort: httpPort,
		handler:    event,
	}
}

// RunServer starts HTTP server
func (r *Router) RunServer(ctx context.Context) {
	engine := gin.Default()
	r.assignRoutes(engine)

	srv := &http.Server{
		Addr:    r.serverPort.HTTPPort,
		Handler: engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

// AssignRoutes assign the available routes of the HTTP API
func (r *Router) assignRoutes(engine *gin.Engine) {
	engine.GET("/bucket", r.handler.ProductsInBucket)
	engine.GET("/outofbucket", r.handler.ProductsOutFromBucket)
	engine.GET("/description", r.handler.ProductsDescription)
	engine.GET("/descriptionandbucket", r.handler.ProductsBucketAndDescription)
}
