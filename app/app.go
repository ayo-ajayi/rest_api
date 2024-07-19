package app

import (
	"context"
	"log"
	"net/http"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ayo-ajayi/rest_api_template/controllers"
	"github.com/ayo-ajayi/rest_api_template/db"
	"github.com/ayo-ajayi/rest_api_template/route"
)

type App struct {
	server   *http.Server
	database *db.DBClient
}

func NewApp() (*App, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	database, err := db.DBinit(ctx)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	ctrl := controllers.NewController(database)
	router, err := route.SetupRouter(ctrl)
	if err != nil {
		log.Fatal(err)
	}
	return &App{
		server: &http.Server{
			Addr:    ":8000",
			Handler: router,
		}, database: database}, nil
}

func (a *App) Start() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed to start: %v", err)
		}
	}()
	log.Println("Server is running🎉🎉. Press Ctrl+C to stop")

	<-stop
	a.Shutdown()
}

func (a *App) Shutdown() {
	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server shutdown completed.")
	if err := a.database.Close(); err != nil {
		log.Printf("Error closing the database: %v", err)
	}
	log.Println("Database successfully closed!!")
}
