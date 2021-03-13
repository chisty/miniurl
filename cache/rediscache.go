package cache

import (
	"encoding/json"
	"log"

	"github.com/chisty/shortlink/model"
	"github.com/gomodule/redigo/redis"
)

type rediscache struct {
	host string
	ttl  int
	pool *redis.Pool
	l    *log.Logger
}

func NewRedis(host string, lg *log.Logger, ttl, mxidle, mxactive int) Cache {
	pl := initPool(host, mxidle, mxactive)
	r := &rediscache{
		host: host,
		ttl:  ttl,
		l:    lg,
		pool: pl,
	}
	return r
}

func (rc *rediscache) Get(key string) (*model.ShortLink, error) {
	conn := rc.pool.Get()
	defer conn.Close()

	val, err := getVal(conn, rc.ttl, key)
	if err != nil {
		rc.l.Printf("redis cache miss: fail to get key %s\n", key)
		return nil, err
	}

	slink := model.ShortLink{}
	err = json.Unmarshal([]byte(val), &slink)
	if err != nil {
		rc.l.Println("redis value unmarshalling error: ", err.Error())
		return nil, err
	}

	rc.l.Println("data found in redis: ", slink.URL)

	return &slink, nil
}

func (rc *rediscache) Set(key string, val *model.ShortLink) error {
	jval, err := json.Marshal(val)
	if err != nil {
		rc.l.Println("redis set error: ", err.Error())
		return err
	}

	conn := rc.pool.Get()
	defer conn.Close()

	err = setVal(conn, rc.ttl, key, string(jval))
	if err != nil {
		rc.l.Printf("redis error: fail to save key %s, error %s\n", key, err.Error())
		return err
	}

	rc.l.Println("data saved in redis: ", string(jval))

	return nil
}

func setVal(conn redis.Conn, ttl int, key, val string) error {
	_, err := conn.Do("SETEX", key, ttl, val)
	if err != nil {
		return err
	}
	return nil
}

func getVal(conn redis.Conn, ttl int, key string) (string, error) {
	val, err := redis.String(conn.Do("GETEX", key, "EX", ttl))
	if err != nil {
		return "", err
	}

	return val, nil
}

func initPool(host string, mxidle, mxactive int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:   mxidle,
		MaxActive: mxactive,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", host)
			if err != nil {
				panic(err)
			}
			return conn, err
		},
	}
}
