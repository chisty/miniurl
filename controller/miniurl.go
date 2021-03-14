package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
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
	Start(response http.ResponseWriter, r *http.Request)
}

func NewMiniURLCtrl(service service.MiniURLSvc, cache cache.Cache, log *log.Logger) MiniURLCtrl {
	return &controller{
		service: service,
		cache:   cache,
		logger:  log,
	}
}

func (ctrl *controller) Start(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Miniurl app is running."))
}

func (ctrl *controller) Get(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := strings.TrimSpace(vars["id"])
	if len(id) == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("invalid request value in route"))
		return
	}

	//try get from cache first
	item, err := ctrl.cache.Get(id)
	if item != nil {
		rw.WriteHeader(http.StatusOK)
		item.ToJSON(rw)
		return
	}

	//try get from db if cache miss
	miniurl, err := ctrl.service.Get(id)
	if err != nil || miniurl == nil {
		rw.WriteHeader(http.StatusNotFound)
		json.NewEncoder(rw).Encode("no value found for this id")
		return
	}

	//set cache with the newly fetched value from db
	ctrl.cache.Set(id, miniurl)

	rw.WriteHeader(http.StatusTemporaryRedirect)
	miniurl.ToJSON(rw)
}

func (ctrl *controller) Save(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	var req model.PostRequest
	err := req.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "could not decode request values", http.StatusBadRequest)
		return
	}

	_, err = url.ParseRequestURI(req.URL)
	if err != nil {
		ctrl.logger.Println(err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	miniurl := model.MiniURL{
		URL:       req.URL,
		CreatedOn: time.Now().UTC().String(),
	}

	item, err := ctrl.service.Save(&miniurl)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode(err)
		return
	}

	//save to cache if data saved in db
	ctrl.cache.Set(miniurl.ID, &miniurl)

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(item)
}
