package config

import "os"

type (
	Config struct {
		HTTPServer HTTPServerConfig
		EventData  EventDataConfig
	}

	HTTPServerConfig struct {
		HTTPServerPort string
	}

	EventDataConfig struct {
		ProjectID string
		TopicID   string
	}
)

func Set() (*Config, error) {
	httpServerPort := os.Getenv("PORT")
	if httpServerPort == "" {
		httpServerPort = ":8080"
	}

	httpServer := HTTPServerConfig{HTTPServerPort: httpServerPort}

	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		projectID = "pub-sub-359008"
	}

	topicID := os.Getenv("TOPIC_ID")
	if topicID == "" {
		topicID = "my-topic"
	}

	eventData := EventDataConfig{
		ProjectID: projectID,
		TopicID:   topicID,
	}

	return &Config{
			HTTPServer: httpServer,
			EventData:  eventData,
		},
		nil
}
