package router

import (
	"Intern/gcp_pub-sub/modules/publisher/pkg/app/config"
	"Intern/gcp_pub-sub/modules/publisher/pkg/handler"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type Router struct {
	httpServerPort *config.HTTPServerConfig
	handler        handler.Handler
}

func New(port *config.HTTPServerConfig, event handler.Handler) *Router {
	return &Router{
		httpServerPort: port,
		handler:        event,
	}
}

func (r *Router) RunServer(ctx context.Context) {
	engine := gin.Default()
	r.assignRoutes(engine)

	srv := &http.Server{
		Addr:    r.httpServerPort.HTTPServerPort,
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

func (r *Router) assignRoutes(engine *gin.Engine) {
	engine.LoadHTMLGlob("./templates/index.html")
	engine.GET("/index", r.handler.Index)
	engine.POST("/send", r.handler.SendData)
}
