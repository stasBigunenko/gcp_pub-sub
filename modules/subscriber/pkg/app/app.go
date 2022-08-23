package app

import (
	"Intern/gcp_pub-sub/modules/subscriber/handler/product"
	"Intern/gcp_pub-sub/modules/subscriber/pkg/app/config"
	"Intern/gcp_pub-sub/modules/subscriber/pkg/app/router"
	"Intern/gcp_pub-sub/modules/subscriber/pkg/app/storage/postgresql"
	"Intern/gcp_pub-sub/modules/subscriber/pkg/app/worker"
	product2 "Intern/gcp_pub-sub/modules/subscriber/repo/product"
	"Intern/gcp_pub-sub/modules/subscriber/service"
	"context"
	"log"
	"os"
)

type Application struct {
	Configuration *config.Configuration
	Worker        *worker.Worker
	Service       service.Service
}

func Create() (*Application, error) {
	cfg, err := config.Set()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	storage, err := postgresql.New(&cfg.StoragePostgreSQL)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if err = storage.AddSomeDataToDB(); err != nil {
		log.Fatalf("error adding data to db: %w\n", err)
	}

	repo := product2.New(storage.Pdb)

	service := product.New(&repo)

	worker := worker.New(repo, &cfg.PubSubData)

	return &Application{
		Configuration: cfg,
		Worker:        worker,
		Service:       service,
	}, nil
}

func (app *Application) Run(ctx context.Context) error {
	app.Worker.Run(ctx)

	router.New(&app.Configuration.HTTPServerPort, product.New(&app.Service)).RunServer(ctx)
	return nil
}
