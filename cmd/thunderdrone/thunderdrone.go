package main

import (
	"context"
	"fmt"
	database2 "github.com/superkruger/thunderdrone/internal/database"
	"github.com/superkruger/thunderdrone/internal/lnd"
	"github.com/superkruger/thunderdrone/internal/server"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	log.Println("Thunderdrone starting...")

	time.Sleep(5 * time.Second)

	db, err := database2.PgConnect("thunderdrone_db", "thunderdrone", "password", "thunderdrone-db", "5432")
	if err != nil {
		log.Println("Could not connect to db", err)
		os.Exit(1)
	}

	err = database2.MigrateUp(db)
	if err != nil {
		log.Println(err)
	}

	var wg sync.WaitGroup
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	go waitForInterruptAndCancel(interrupt, cancel)

	wg.Add(1)
	go func(context context.Context) {
		server.Start(context, db, "8080")
		log.Println("8080 stopped")
		wg.Done()
	}(ctx)

	wg.Add(1)
	go func(context context.Context) {
		lnd.Start(context, db)
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

func waitForInterruptAndCancel(interrupt chan os.Signal, cancel context.CancelFunc) {
	<-interrupt
	log.Println("Interrupt received, cancelling context")
	cancel()
}
