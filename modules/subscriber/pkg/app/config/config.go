package config

import (
	"github.com/google/uuid"
	"os"
)

type (
	// Configuration object is the main object
	Configuration struct {
		HTTPServerPort    HTTPServerConfiguration
		PubSubData        PubSubConfiguration
		StoragePostgreSQL StorageConfiguration
	}

	// HTTPServerConfiguration is the object that define http port configuration
	HTTPServerConfiguration struct {
		HTTPPort string
	}

	// PubSubConfiguration is the object that define configuration of the google pub-sub API config
	PubSubConfiguration struct {
		ProjectID    string
		TopicID      string
		SubscriberID string
	}

	// StorageConfiguration is the object that define db connection
	StorageConfiguration struct {
		ConnString string
	}
)

// Set initialize environment variables if they didn't set
func Set() (*Configuration, error) {
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = ":8081"
	}

	httpServer := HTTPServerConfiguration{httpPort}

	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		projectID = "pub-sub-359008"
	}

	topicID := os.Getenv("TOPIC_ID")
	if topicID == "" {
		topicID = "my-topic"
	}

	id := uuid.New().String()
	subscriberID := string([]rune(id)[:0]) + "s" + string([]rune(id)[1:])

	pubSub := PubSubConfiguration{
		ProjectID:    projectID,
		TopicID:      topicID,
		SubscriberID: subscriberID,
	}

	hostDB := os.Getenv("HOST_DB")
	if hostDB == "" {
		hostDB = "localhost"
	}

	portDB := os.Getenv("PORT_DB")
	if portDB == "" {
		portDB = "5432"
	}

	userDB := os.Getenv("USER_DB")
	if userDB == "" {
		userDB = "pub-sub"
	}

	pswDB := os.Getenv("PSW_DB")
	if pswDB == "" {
		pswDB = "qwerty"
	}

	nameDB := os.Getenv("NAME_DB")
	if nameDB == "" {
		nameDB = "pub-sub"
	}

	ssldb := os.Getenv("SSLDB")
	if ssldb == "" {
		ssldb = "disable"
	}

	connStr := "host=" + hostDB + " port=" + portDB + " user=" + userDB + " password=" + pswDB + " dbname=" + nameDB + " sslmode=" + ssldb

	storage := StorageConfiguration{connStr}

	return &Configuration{
			HTTPServerPort:    httpServer,
			PubSubData:        pubSub,
			StoragePostgreSQL: storage,
		},
		nil
}
