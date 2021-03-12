package service

import (
	"log"

	"github.com/chisty/shortlink/database"
	"github.com/chisty/shortlink/model"
	"github.com/google/uuid"
)

//ShortLinkService ---
type LinkService interface {
	Get(id string) (*model.ShortLink, error)
	Save(data *model.ShortLink) (*model.ShortLink, error)
}

type service struct {
	db database.DB
	l  *log.Logger
}

//NewService ---
func NewService(d database.DB, log *log.Logger) LinkService {
	return &service{
		db: d,
		l:  log,
	}
}

func (s *service) Get(id string) (*model.ShortLink, error) {
	s.l.Printf("service get %s", id)
	return s.db.Get(id)
}

func (s *service) Save(data *model.ShortLink) (*model.ShortLink, error) {
	data.ID = getShaValue(data.URL)

	s.l.Printf("service save %s -> %s", data.URL, data.ID)

	err := s.db.Save(data)
	return data, err
}

//save sha value for now.
func getShaValue(text string) string {
	id := uuid.New()
	return id.String()
}
