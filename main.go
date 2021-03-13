package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chisty/miniurl/cache"
	"github.com/chisty/miniurl/controller"
	"github.com/chisty/miniurl/database"
	"github.com/chisty/miniurl/service"
	"github.com/gorilla/mux"
)

func main() {
	cfg := LoadConfig()
	fmt.Printf("%+v", cfg)

	logger := log.New(os.Stdout, "shortlink-app", log.LstdFlags|log.Lshortfile)
	redis := cache.NewRedis(cfg.Redis.Host, logger, cfg.Redis.TTL, cfg.Redis.MaxIdle, cfg.Redis.MaxActive)
	db := database.NewDynamoDB(cfg.AWS.Table, logger, cfg.AWS.Region, cfg.AWS.AccessKey, cfg.AWS.SecretKey)
	svc := service.NewMiniURLSvc(db, cfg.IdGenNodeId, logger)
	ctrl := controller.NewMiniURLCtrl(svc, redis, logger)

	router := mux.NewRouter()
	getRouter := router.Methods(http.MethodGet).Subrouter()
	postRouter := router.Methods(http.MethodPost).Subrouter()

	// getRouter.HandleFunc("/test", ctrl.Test)
	getRouter.HandleFunc("/{id}", ctrl.Get)
	postRouter.HandleFunc("/", ctrl.Save)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
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

	logger.Println("server running...")

	select {
	case err := <-appFatalError:
		if err != nil {
			logger.Fatal(err)
			return
		}
	case sig := <-sigChan:
		logger.Println("received termination command. Shutting down gracefully. signal= ", sig)
		tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
		err := s.Shutdown(tc)
		if err != nil {
			logger.Printf("main: graceful shutdown failed %v", err)
			return
		}
	}
}
