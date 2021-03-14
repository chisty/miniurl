package service

import (
	"log"

	"github.com/bwmarrin/snowflake"
	"github.com/chisty/miniurl/database"
	"github.com/chisty/miniurl/model"
)

var base64 []rune

func init() {
	//base64 is the selected 64 unique digits which is used for base conversion
	base64 = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_=")
}

type service struct {
	db     database.DB
	idGen  *snowflake.Node
	logger *log.Logger
}

//NewMiniURLSvc is the implementation of MiniURLSvc
func NewMiniURLSvc(d database.DB, idgenId int, log *log.Logger) MiniURLSvc {
	//idGen is the id generator. Currently it is taking idGenID as node id. In real scenario, this should be a separate service which
	//will provide unique id for a distributed system. snowflake is twitter's approach to create unique id for distributed system.
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

func (svc *service) Get(id string) (*model.MiniURL, error) {
	return svc.db.Get(id)
}

func (svc *service) Save(data *model.MiniURL) (*model.MiniURL, error) {
	data.ID = nextID(svc.idGen)
	err := svc.db.Save(data)
	return data, err
}

func nextID(idGen *snowflake.Node) string {
	id := idGen.Generate().Int64()
	return convertToBase64(id, base64)
}

//basic base conversion from decimal to base64
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
