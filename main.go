package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/chisty/miniurl/cache"
	"github.com/chisty/miniurl/controller"
	"github.com/chisty/miniurl/database"
	"github.com/chisty/miniurl/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

var mySigningKey = []byte("my_super_secret_key")

func handleAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("Middleware Auth")
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			fmt.Println("Malformed token")
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Malformed token."))
		} else {
			jwtToken := authHeader[1]
			token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(mySigningKey), nil
			})

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				fmt.Println(claims)
				fmt.Println("----------------------")
				fmt.Printf("%+v\n", claims)
				ctx := context.WithValue(r.Context(), "props", claims)
				next.ServeHTTP(rw, r.WithContext(ctx))
			} else {
				fmt.Println(err)
				rw.WriteHeader(http.StatusUnauthorized)
				rw.Write([]byte("Unauthorized"))
			}
		}
	})
}

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
	getRouter.HandleFunc("/{id}", handleAuth(ctrl.Get))
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
