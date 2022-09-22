package app

import (
	"Intern/gcp_pub-sub/modules/subscriber/handler/product"
	"Intern/gcp_pub-sub/modules/subscriber/pkg/app/config"
	"Intern/gcp_pub-sub/modules/subscriber/pkg/app/router"
	"Intern/gcp_pub-sub/modules/subscriber/pkg/app/storage/postgresql"
	"Intern/gcp_pub-sub/modules/subscriber/pkg/app/worker"
	repo "Intern/gcp_pub-sub/modules/subscriber/repo/product"
	"Intern/gcp_pub-sub/modules/subscriber/service"
	serv "Intern/gcp_pub-sub/modules/subscriber/service/product"
	"context"
	"log"
	"os"
)

// Application Config is the top-level configuration object.
type Application struct {
	Configuration *config.Configuration
	Worker        *worker.Worker
	Service       service.Service
}

// Create is a constructor (initialize) of the Application object
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

	if err = storage.AddSomeDataIntoTable(cfg.StoragePostgreSQL.Table1.Name, cfg.StoragePostgreSQL.Table1.Path, cfg.StoragePostgreSQL.Table1.FieldsQty); err != nil {
		log.Fatalf("error adding data to db: %v\n", err)
	}

	if err = storage.AddSomeDataIntoTable(cfg.StoragePostgreSQL.Table2.Name, cfg.StoragePostgreSQL.Table2.Path, cfg.StoragePostgreSQL.Table2.FieldsQty); err != nil {
		log.Fatalf("error adding data to db: %v\n", err)
	}

	if err = storage.AddSomeDataIntoTable(cfg.StoragePostgreSQL.Table3.Name, cfg.StoragePostgreSQL.Table3.Path, cfg.StoragePostgreSQL.Table3.FieldsQty); err != nil {
		log.Fatalf("error adding data to db: %v\n", err)
	}

	repo := repo.New(storage.Pdb)

	service := serv.New(repo)

	worker := worker.New(repo, &cfg.PubSubData)

	return &Application{
		Configuration: cfg,
		Worker:        worker,
		Service:       service,
	}, nil
}

// Run is the function that start application
func (app *Application) Run(ctx context.Context) error {
	app.Worker.Run(ctx)

	router.New(&app.Configuration.HTTPServerPort, product.New(app.Service)).RunServer(ctx)
	return nil
}
