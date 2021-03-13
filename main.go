package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chisty/shortlink/cache"
	"github.com/chisty/shortlink/controller"
	"github.com/chisty/shortlink/database"
	"github.com/chisty/shortlink/service"
	"github.com/gorilla/mux"
)

var (
	table string = os.Getenv("AWS_TABLE")
	nid   string = os.Getenv("NODE_ID")

	l     *log.Logger               = log.New(os.Stdout, "shortlink-app", log.LstdFlags|log.Lshortfile)
	redis cache.Cache               = cache.NewRedis("localhost:6379", l, 10, 80, 12000)
	db    database.DB               = database.NewDynamoDB(table, l)
	svc   service.LinkService       = service.NewService(db, nid, l)
	ctrl  controller.LinkController = controller.NewLinkController(svc, redis, l)
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

	appFatalError := make(chan error, 1)

	go func() {
		appFatalError <- s.ListenAndServe()
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	l.Println("server running...")

	select {
	case err := <-appFatalError:
		if err != nil {
			l.Fatal(err)
			return
		}
	case sig := <-sigChan:
		l.Println("received termination command. Shutting down gracefully. signal= ", sig)
		tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
		err := s.Shutdown(tc)
		if err != nil {
			l.Printf("main: graceful shutdown failed %v", err)
			return
		}
	}
}
