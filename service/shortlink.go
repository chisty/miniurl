package service

import (
	"log"

	"github.com/bwmarrin/snowflake"
	"github.com/chisty/shortlink/database"
	"github.com/chisty/shortlink/model"
)

var base64 []rune

//ShortLinkService ---
type LinkService interface {
	Get(id string) (*model.ShortLink, error)
	Save(data *model.ShortLink) (*model.ShortLink, error)
}

type service struct {
	db     database.DB
	idGen  *snowflake.Node
	logger *log.Logger
}

func init() {
	base64 = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_=")
}

//NewService ---
func NewService(d database.DB, idgenId int, log *log.Logger) LinkService {
	idGen, err := snowflake.NewNode(int64(idgenId))
	if err != nil {
		panic(err)
	}

	return &service{
		db:     d,
		idGen:  idGen,
		logger: log,
	}
}

func (s *service) Get(id string) (*model.ShortLink, error) {
	s.logger.Printf("service get %s", id)
	return s.db.Get(id)
}

func (s *service) Save(data *model.ShortLink) (*model.ShortLink, error) {
	data.ID = getNextID(s.idGen)
	err := s.db.Save(data)
	return data, err
}

func getNextID(idGen *snowflake.Node) string {
	id := idGen.Generate().Int64()
	return convertToBase64(id, base64)
}

func convertToBase64(val int64, baseD []rune) string {
	temp := []rune{}
	baseL := int64(len(baseD))
	neg := false

	if val == 0 {
		return "0"
	}

	if val < 0 {
		neg = true
		val *= -1
	}

	for val > 0 {
		temp = append(temp, baseD[val%baseL])
		val /= baseL
	}

	l := len(temp)
	for i := 0; i < l/2; i++ {
		temp[i], temp[l-i-1] = temp[l-i-1], temp[i]
	}

	if neg {
		return "-" + string(temp)
	}

	return string(temp)
}
