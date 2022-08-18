package config

import "os"

type Config struct {
	Port      string
	ProjectID string
	TopicID   string
}

func Set() (*Config, error) {
	var config Config

	config.Port = os.Getenv("PORT")
	if config.Port == "" {
		config.Port = ":8080"
	}

	config.ProjectID = os.Getenv("PROJECT_ID")
	if config.ProjectID == "" {
		config.ProjectID = "pub-sub-359008"
	}

	config.TopicID = os.Getenv("TOPIC_ID")
	if config.TopicID == "" {
		config.TopicID = "my-topic"
	}

	return &Config{
			Port:      config.Port,
			ProjectID: config.ProjectID,
			TopicID:   config.TopicID,
		},
		nil
}
