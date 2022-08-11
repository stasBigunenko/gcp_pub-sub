package handler

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"gcp_pub-sub/cmd/front-end_publisher/config"
	"gcp_pub-sub/internal/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func init() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "C:\\goProjects\\src\\Intern\\gcp_pub-sub\\pub-sub-359008-ff94c59da4aa.json")
}

type Handler struct {
	Config *config.Config
}

func New(cfg *config.Config) *Handler {
	return &Handler{
		Config: cfg,
	}
}

func (h *Handler) Routes(r *gin.Engine) *gin.Engine {
	r.LoadHTMLGlob("./cmd/front-end_publisher/templates/*.html")
	r.GET("/index", h.Index)
	r.POST("/send", h.SendData)

	return r
}

func (h *Handler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"topic": h.Config.TopicID,
	})
}

func (h *Handler) SendData(c *gin.Context) {
	var message models.Message

	message.CategoryID = c.PostForm("categoryID")
	message.ProductID = c.PostForm("productID")
	message.ActionID = c.PostForm("actionID")

	data, err := json.Marshal(message)
	if err != nil {
		c.Error(err)
		return
	}

	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, h.Config.ProjectID)
	if err != nil {
		c.Error(fmt.Errorf("Could not create pubsub Client: %v", err))
		return
	}

	if err = publish(client, h.Config.TopicID, data); err != nil {
		c.Error(fmt.Errorf("Failed to publish: %v", err))
		return
	}

	c.Redirect(301, "/index")
}

func publish(client *pubsub.Client, topicID string, msg []byte) error {
	ctx := context.Background()
	topic := client.Topic(topicID)

	result := topic.Publish(ctx, &pubsub.Message{
		Data: msg,
	})

	id, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("pubsub: result.Get: %v", err)
	}

	log.Printf("Published a message; msg ID: %v\n", id)

	return nil
}
