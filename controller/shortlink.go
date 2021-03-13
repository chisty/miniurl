package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/chisty/shortlink/cache"
	"github.com/chisty/shortlink/model"
	"github.com/chisty/shortlink/service"
	"github.com/gorilla/mux"
)

type linkController struct {
	service service.LinkService
	cache   cache.Cache
	logger  *log.Logger
}

type LinkController interface {
	Get(response http.ResponseWriter, r *http.Request)
	Save(response http.ResponseWriter, r *http.Request)
}

func NewLinkController(service service.LinkService, cache cache.Cache, log *log.Logger) LinkController {
	return &linkController{
		service: service,
		cache:   cache,
		logger:  log,
	}
}

func (lc *linkController) Get(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	item, err := lc.cache.Get(id)
	if item != nil {
		rw.WriteHeader(http.StatusOK)
		item.ToJSON(rw)
		return
	}

	slink, err := lc.service.Get(id)
	if err != nil || slink == nil {
		rw.WriteHeader(http.StatusNotFound)
		json.NewEncoder(rw).Encode("no value found for this id")
		return
	}

	lc.cache.Set(id, slink)

	rw.WriteHeader(http.StatusOK)
	slink.ToJSON(rw)
}

func (lc *linkController) Save(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	link := model.ShortLink{}
	err := link.FromJSON(r.Body)

	if err != nil {
		http.Error(rw, "invalid or malformed input", http.StatusBadRequest)
		return
	}

	link.CreatedOn = time.Now().UTC().String()

	item, err := lc.service.Save(&link)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode(err)
		return
	}

	lc.cache.Set(link.ID, &link)

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(item)
}
