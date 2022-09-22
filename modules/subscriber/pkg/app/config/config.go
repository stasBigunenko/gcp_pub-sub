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
		Table1     Table
		Table2     Table
		Table3     Table
	}

	Table struct {
		Name      string
		Path      string
		FieldsQty int
	}
)

const (
	TABLE1_NAME = "products"
	PATH1       = "./data/products.csv"
	FIELDS1Qty  = 5

	TABLE2_NAME = "categories"
	PATH2       = "./data/categories.csv"
	FIELDS2Qty  = 2

	TABLE3_NAME = "actions"
	PATH3       = "./data/categories.csv"
	FIELDS3Qty  = 2
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
		projectID = "pub-sub-46957"
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

	table1 := Table{TABLE1_NAME, PATH1, FIELDS1Qty}
	table2 := Table{TABLE2_NAME, PATH2, FIELDS2Qty}
	table3 := Table{TABLE3_NAME, PATH3, FIELDS3Qty}

	storage := StorageConfiguration{connStr, table1, table2, table3}

	return &Configuration{
			HTTPServerPort:    httpServer,
			PubSubData:        pubSub,
			StoragePostgreSQL: storage,
		},
		nil
}
