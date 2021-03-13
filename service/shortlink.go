package service

import (
	"fmt"
	"log"
	"strconv"

	"github.com/bwmarrin/snowflake"
	"github.com/chisty/shortlink/database"
	"github.com/chisty/shortlink/model"
)

var base64 []rune

//ShortLinkService ---
type LinkService interface {
	Get(id string) (*model.ShortLink, error)
	Save(data *model.ShortLink) (*model.ShortLink, error)
	Test()
}

type service struct {
	db   database.DB
	node *snowflake.Node
	l    *log.Logger
}

func init() {
	base64 = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_=")
}

//NewService ---
func NewService(d database.DB, nid string, log *log.Logger) LinkService {
	n := getSnowflakeNode(nid)
	return &service{
		db:   d,
		node: n,
		l:    log,
	}
}

func (s *service) Test() {
	fmt.Println("Test done")
}

func (s *service) Get(id string) (*model.ShortLink, error) {
	s.l.Printf("service get %s", id)
	return s.db.Get(id)
}

func (s *service) Save(data *model.ShortLink) (*model.ShortLink, error) {
	data.ID = getNextID(s.node)
	err := s.db.Save(data)
	return data, err
}

func getNextID(node *snowflake.Node) string {
	id := node.Generate().Int64()
	return convertToBase64(id, base64)
}

func convertToBase64(val int64, baseD []rune) string {
	temp := []rune{}
	baseL := int64(len(baseD))

	for val > 0 {
		temp = append(temp, baseD[val%baseL])
		val /= baseL
	}

	return string(temp)
}

func getSnowflakeNode(nid string) *snowflake.Node {
	if len(nid) == 0 {
		nid = "1"
	}
	n, err := strconv.ParseInt(nid, 10, 64)
	if err != nil {
		panic(err)
	}
	node, err := snowflake.NewNode(n)
	if err != nil {
		panic(err)
	}

	return node
}
