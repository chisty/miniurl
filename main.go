package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/chisty/shortlink/controller"
	"github.com/chisty/shortlink/database"
	"github.com/chisty/shortlink/service"
	"github.com/gorilla/mux"
)

var (
	table string                    = os.Getenv("AWS_TABLE")
	nid   string                    = os.Getenv("NODE_ID")
	l     *log.Logger               = log.New(os.Stdout, "shortlink-app", log.LstdFlags|log.Lshortfile)
	db    database.DB               = database.NewDynamoDB(table, l)
	svc   service.LinkService       = service.NewService(db, nid, l)
	ctrl  controller.LinkController = controller.NewLinkController(svc, l)
)

func main() {
	router := mux.NewRouter()
	getRouter := router.Methods(http.MethodGet).Subrouter()
	postRouter := router.Methods(http.MethodPost).Subrouter()

	getRouter.HandleFunc("/test", ctrl.Test)
	getRouter.HandleFunc("/{id}", ctrl.Get)
	postRouter.HandleFunc("/", ctrl.Save)

	s := &http.Server{
		Addr:         ":9000",
		Handler:      router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	l.Println("Server running...")

	sig := <-sigChan
	l.Println("Received termination command. Shutting down gracefully. signal= ", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
