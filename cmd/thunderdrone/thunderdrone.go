package main

import (
	"context"
	"fmt"
	"github.com/superkruger/thunderdrone/internal/database"
	"github.com/superkruger/thunderdrone/internal/lnd"
	"github.com/superkruger/thunderdrone/internal/repositories"
	"github.com/superkruger/thunderdrone/internal/routes"
	"github.com/superkruger/thunderdrone/internal/server"
	"github.com/superkruger/thunderdrone/internal/services"
	"log"
	"os"
	"os/signal"
	"sync"
)

func main() {
	log.Println("Thunderdrone starting...")

	db, err := database.PgConnect("thunderdrone_db", "thunderdrone", "password", "thunderdrone-db", "5432")
	if err != nil {
		log.Println("Could not connect to db", err)
		os.Exit(1)
	}

	err = database.MigrateUp(db)
	if err != nil {
		log.Println(err)
	}

	var wg sync.WaitGroup
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	go waitForInterruptAndCancel(interrupt, cancel)

	nodeSettings := services.NewNodeSettingsService(repositories.NewNodeSettingsRepo(db))

	wg.Add(1)
	go func(context context.Context) {

		routes := allRoutes(nodeSettings)

		server.Start(context, routes, "8080")
		log.Println("8080 stopped")
		wg.Done()
	}(ctx)

	wg.Add(1)
	go func(context context.Context) {

		lndClient := lnd.NewLndClient(context, nodeSettings)

		lndClient.Start()
		fmt.Println("lnd done.")
		wg.Done()
	}(ctx)

	wg.Wait()

	log.Println("Closing db")
	err = db.Close()
	if err != nil {
		log.Println("Error closing db", err)
	}
	log.Println("All done")
}

func allRoutes(nodeSettings services.NodeSettingsService) []routes.Routable {
	return []routes.Routable{
		routes.NewNodeSettingsRoutes(nodeSettings),
	}
}

func waitForInterruptAndCancel(interrupt chan os.Signal, cancel context.CancelFunc) {
	<-interrupt
	log.Println("Interrupt received, cancelling context")
	cancel()
}
