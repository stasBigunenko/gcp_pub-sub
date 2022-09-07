package worker

import (
	"Intern/gcp_pub-sub/modules/subscriber/model"
	"Intern/gcp_pub-sub/modules/subscriber/pkg/app/config"
	"Intern/gcp_pub-sub/modules/subscriber/repo"
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// Worker is the object that define entity which works with Google API and Miroservice API
type Worker struct {
	client       *pubsub.Client
	repo         repo.ProductsRepo
	pubSubConfig *config.PubSubConfiguration
}

// init initialize Goggle app environment variable
func init() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./pub-sub-46957-54560624752f.json")
}

// New is the contructor of the Worker entity
func New(r repo.ProductsRepo, c *config.PubSubConfiguration) *Worker {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, c.ProjectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	topic := createTopicIfNotExists(client, c.TopicID)
	if err != nil {
		log.Fatalf("Failed to create topic: %v", err)
	}

	if err = create(client, c.SubscriberID, topic); err != nil {
		log.Fatal(err)
	}

	return &Worker{
		client:       client,
		repo:         r,
		pubSubConfig: c,
	}
}

// Run makes our worker work
func (w *Worker) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var err error

	go func() {
		if err = w.pullMsgs(w.client, w.pubSubConfig.SubscriberID); err != nil {
			cancel()
			return
		}
	}()

	if err != nil {
		return err
	}

	return nil
}

// GetCLient return pub-sub client connection
func (w *Worker) GetClient() *pubsub.Client {
	return w.client
}

// PullMsgs pull messages from Google Pub-Sub
func (w *Worker) pullMsgs(client *pubsub.Client, name string) error {
	ctx := context.Background()

	var mu sync.Mutex
	sub := client.Subscription(name)

	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()

		var body model.Message

		if err := json.Unmarshal(msg.Data, &body); err != nil {
			log.Printf("error in Unmarshal: %v\n", err)
			return
		}

		w.repo.AddAction(body.ActionID, body.ProductID)

		mu.Lock()
		defer mu.Unlock()
	})

	if err != nil {
		return err
	}

	return nil
}

// Create subscription in the Google Pub-Sub
func create(client *pubsub.Client, name string, topic *pubsub.Topic) error {
	ctx := context.Background()

	sub, err := client.CreateSubscription(ctx, name, pubsub.SubscriptionConfig{
		Topic:       topic,
		AckDeadline: 20 * time.Second,
	})
	if err != nil {
		log.Printf("err in subscription: %v\n", err)
		return nil
	}
	fmt.Printf("Created subscription: %v\n", sub)

	return nil
}

// CreateTopicIfNotExists creates topic in the Google Pub-Sub
func createTopicIfNotExists(c *pubsub.Client, topic string) *pubsub.Topic {
	ctx := context.Background()
	t := c.Topic(topic)
	ok, err := t.Exists(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if ok {
		return t
	}
	t, err = c.CreateTopic(ctx, topic)
	if err != nil {
		log.Fatalf("Failed to create the topic: %v", err)
	}
	return t
}
