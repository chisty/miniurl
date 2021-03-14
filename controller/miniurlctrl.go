package controller

import "net/http"

//MiniURLCtrl is controller interface
type MiniURLCtrl interface {
	Get(response http.ResponseWriter, r *http.Request)
	Save(response http.ResponseWriter, r *http.Request)
	Start(response http.ResponseWriter, r *http.Request)
}
