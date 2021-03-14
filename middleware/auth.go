package Middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type Auth struct {
	secretKey []byte
	logger    *log.Logger
}

func NewAuth(sk string, log *log.Logger) *Auth {
	return &Auth{
		secretKey: []byte(sk),
		logger:    log,
	}
}

func (auth *Auth) HandleAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized. Malformed token."))
		} else {
			jwtToken := authHeader[1]
			token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					auth.logger.Fatalf("unexpected signing method: %v", token.Header["alg"])
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(auth.secretKey), nil
			})

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				if appid, ok := claims["app"]; ok {
					if strings.ToLower(appid.(string)) == "miniurl" {
						ctx := context.WithValue(r.Context(), "props", claims)
						next.ServeHTTP(rw, r.WithContext(ctx))
						return
					}
				}

			}

			auth.logger.Println(err)
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Unauthorized"))
		}
	})
}
