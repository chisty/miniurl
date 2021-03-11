package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

type LongURL struct {
	URL string `json:"URL"`
}

func main() {
	router := mux.NewRouter()
	getRouter := router.Methods(http.MethodGet).Subrouter()
	postRouter := router.Methods(http.MethodPost).Subrouter()

	getRouter.HandleFunc("/{id}", getByID)

	postRouter.HandleFunc("/", save)

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

	fmt.Println("Server running...")

	sig := <-sigChan
	fmt.Println("Received termination command. Shutting down gracefully. signal= ", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}

func getByID(response http.ResponseWriter, r *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println("ShortUrl: ", id)

	response.WriteHeader(http.StatusOK)
}

func save(response http.ResponseWriter, r *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	input := LongURL{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)

	if err != nil {
		panic(err)
	}
	fmt.Println(input.URL)
	response.WriteHeader(http.StatusOK)
}
