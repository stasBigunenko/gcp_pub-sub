package event

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"gcp_pub-sub/modules/publisher/pkg/app/config"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func init() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "C:\\goProjects\\src\\Intern\\gcp_pub-sub\\pub-sub-359008-ff94c59da4aa.json")
}

type Event struct {
	config *config.Config
}

type message struct {
	ProductID string `json:"productID"`
	ActionID  string `json:"actionID"`
}

func New(cfg *config.Config) *Event {
	return &Event{
		config: cfg,
	}
}

func (e *Event) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"topic": e.config.TopicID,
	})
}

func (e *Event) SendData(c *gin.Context) {
	var msg message

	msg.ProductID = c.PostForm("productID")
	msg.ActionID = c.PostForm("actionID")

	data, err := json.Marshal(msg)
	if err != nil {
		c.Error(err)
		return
	}

	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, e.config.ProjectID)
	if err != nil {
		c.Error(fmt.Errorf("Could not create pubsub Client: %v", err))
		return
	}

	if err = publish(client, e.config.TopicID, data); err != nil {
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
