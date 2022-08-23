package app

import (
	"Intern/gcp_pub-sub/modules/publisher/pkg/app/config"
	"Intern/gcp_pub-sub/modules/publisher/pkg/app/router"
	"Intern/gcp_pub-sub/modules/publisher/pkg/handler/event"
	"context"
	"log"
	"os"
)

type Application struct {
	Configuration *config.Config
}

func Create() (*Application, error) {
	cfg, err := config.Set()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	return &Application{
		Configuration: cfg,
	}, nil
}

func (app *Application) Run(ctx context.Context) error {
	router.New(app.Configuration, event.New(app.Configuration)).RunServer(ctx)
	return nil
}
