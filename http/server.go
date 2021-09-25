package http

import (
	"fmt"
	"log"

	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"

	"github.com/MihaiBlebea/go-event-bus/bus"
	"github.com/MihaiBlebea/go-event-bus/project"
)

const prefix = "/api/v1/"

func New(b bus.Service, proj project.Service, logger *logrus.Logger) {
	r := mux.NewRouter()

	api := r.PathPrefix(prefix).Subrouter()

	// Handle api calls
	api.Handle("/health-check", healthHandler()).
		Methods(http.MethodGet)

	api.Handle("/project", project.CreateHandler(proj)).
		Methods(http.MethodPost)

	api.Handle("/projects", project.ListHandler(proj, b)).
		Methods(http.MethodGet)

	api.Handle("/events", bus.SentEventsHandler(b, proj)).
		Methods(http.MethodGet)

	api.Handle("/subscribe", bus.SubscribeHandler(b, proj)).
		Methods(http.MethodPost)

	api.Handle("/handle", bus.TriggerHandler(b, proj)).
		Methods(http.MethodPost)

	r.Use(loggerMiddleware(logger))

	srv := &http.Server{
		Handler:      cors.AllowAll().Handler(r),
		Addr:         fmt.Sprintf("0.0.0.0:%s", os.Getenv("HTTP_PORT")),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info(fmt.Sprintf("Started server on port %s", os.Getenv("HTTP_PORT")))

	log.Fatal(srv.ListenAndServe())
}
