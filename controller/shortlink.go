package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/chisty/shortlink/model"
	"github.com/chisty/shortlink/service"
	"github.com/gorilla/mux"
)

type linkController struct {
	s service.LinkService
	l *log.Logger
}

type LinkController interface {
	Test(rw http.ResponseWriter, r *http.Request)
	Get(response http.ResponseWriter, r *http.Request)
	Save(response http.ResponseWriter, r *http.Request)
}

func NewLinkController(service service.LinkService, log *log.Logger) LinkController {
	return &linkController{
		s: service,
		l: log,
	}
}

func (lc *linkController) Test(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Test working. Hello from server..."))
}

func (lc *linkController) Get(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]
	lc.l.Println("request id: ", id)

	slink, err := lc.s.Get(id)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		json.NewEncoder(rw).Encode("no data found")
		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(slink)
}

func (lc *linkController) Save(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	link := model.ShortLink{}
	err := link.FromJSON(r.Body)

	if err != nil {
		http.Error(rw, "Unable to unmarshal JSON", http.StatusBadRequest)
		return
	}

	link.CreatedOn = time.Now().UTC().String()

	item, err := lc.s.Save(&link)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode(err)
		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(item)
}
