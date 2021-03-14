package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/chisty/miniurl/cache"
	"github.com/chisty/miniurl/model"
	"github.com/chisty/miniurl/service"
	"github.com/gorilla/mux"
)

type controller struct {
	service service.MiniURLSvc
	cache   cache.Cache
	logger  *log.Logger
}

type MiniURLCtrl interface {
	Get(response http.ResponseWriter, r *http.Request)
	Save(response http.ResponseWriter, r *http.Request)
}

func NewMiniURLCtrl(service service.MiniURLSvc, cache cache.Cache, log *log.Logger) MiniURLCtrl {
	return &controller{
		service: service,
		cache:   cache,
		logger:  log,
	}
}

func (ctrl *controller) Get(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	item, err := ctrl.cache.Get(id)
	if item != nil {
		rw.WriteHeader(http.StatusOK)
		item.ToJSON(rw)
		return
	}

	miniurl, err := ctrl.service.Get(id)
	if err != nil || miniurl == nil {
		rw.WriteHeader(http.StatusNotFound)
		json.NewEncoder(rw).Encode("no value found for this id")
		return
	}

	ctrl.cache.Set(id, miniurl)

	rw.WriteHeader(http.StatusOK)
	miniurl.ToJSON(rw)
}

func (ctrl *controller) Save(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	miniurl := model.MiniURL{}
	err := miniurl.FromJSON(r.Body)

	if err != nil {
		http.Error(rw, "invalid or malformed input", http.StatusBadRequest)
		return
	}

	miniurl.CreatedOn = time.Now().UTC().String()

	item, err := ctrl.service.Save(&miniurl)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode(err)
		return
	}

	ctrl.cache.Set(miniurl.ID, &miniurl)

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(item)
}
