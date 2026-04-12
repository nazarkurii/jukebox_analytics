package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	infra "github.com/nazarkurii/jukebox_analytics/infrastructure"
	"github.com/nazarkurii/jukebox_analytics/internal/shared/httputil"

	"github.com/nazarkurii/jukebox_analytics/internal/handler"
	postresrepo "github.com/nazarkurii/jukebox_analytics/internal/repo/postgres"
	"github.com/nazarkurii/jukebox_analytics/internal/service"
)

type app struct {
	log    *log.Logger
	server *http.Server
}

func (app *app) run() {
	app.log.Println("Analytics server is running...")
	app.server.ListenAndServe()
}

func (app *app) watchForShutdown(ctx context.Context) {
	<-ctx.Done()
	app.log.Print("Analytics server is being shut down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := app.server.Shutdown(ctx)
	if err != nil {
		app.log.Print("Analytics server has been forcibly shut down")
	} else {
		app.log.Print("Analytics server has been shut down")
	}
}

func wireApp(ctx context.Context) (*app, error) {
	server, mux := infra.NewServer()
	app := &app{
		log:    log.Default(),
		server: server,
	}

	db, err := infra.NewPostgres(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	//Logger middleware
	handlerLogger := httputil.NewLogger(app.log)

	//Healthcheck
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	//Playback logger endpoints
	handler.NewPlaybackLogger(service.NewPlaybackLogger(postresrepo.NewLogger(db))).RegisterRoutes(handlerLogger, mux)
	//Stats endpoints
	handler.NewStats(service.NewStats(postresrepo.NewStats(db))).RegisterRoutes(handlerLogger, mux)
	//Track endpoints
	handler.NewTrack(service.NewTrack(postresrepo.NewTrack(db))).RegisterRoutes(handlerLogger, mux)

	return app, nil
}
