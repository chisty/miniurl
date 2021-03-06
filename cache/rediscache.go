package cache

import (
	"encoding/json"
	"log"

	"github.com/chisty/miniurl/model"
	"github.com/gomodule/redigo/redis"
)

type rediscache struct {
	host   string
	ttl    int
	pool   *redis.Pool
	logger *log.Logger
}

//NewRedis is implemenation of redis cache maintaining Cahe interface.
func NewRedis(host string, lg *log.Logger, ttl, mxidle, mxactive int) Cache {
	pl := initPool(host, mxidle, mxactive)
	r := &rediscache{
		host:   host,
		ttl:    ttl,
		logger: lg,
		pool:   pl,
	}
	return r
}

func (rc *rediscache) Get(key string) (*model.MiniURL, error) {
	conn := rc.pool.Get()
	defer conn.Close()

	val, err := getVal(conn, rc.ttl, key)
	if err != nil {
		rc.logger.Printf("redis cache miss: fail to get key %s\n", key)
		return nil, err
	}

	miniurl := model.MiniURL{}
	err = json.Unmarshal([]byte(val), &miniurl)
	if err != nil {
		rc.logger.Println("redis value unmarshalling error: ", err.Error())
		return nil, err
	}

	rc.logger.Println("data found in redis: ", miniurl.URL)

	return &miniurl, nil
}

func (rc *rediscache) Set(key string, val *model.MiniURL) error {
	jval, err := json.Marshal(val)
	if err != nil {
		rc.logger.Println("redis set error: ", err.Error())
		return err
	}

	conn := rc.pool.Get()
	defer conn.Close()

	err = setVal(conn, rc.ttl, key, string(jval))
	if err != nil {
		rc.logger.Printf("redis error: fail to save key %s, error %s\n", key, err.Error())
		return err
	}

	rc.logger.Println("data saved in redis: ", string(jval))

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
	val, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", err
	}

	//if getvalue success try to update ttl
	_ = setVal(conn, ttl, key, val)
	return val, nil
}

//initpool will create redis connection pool which will be faster and easier to access.
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
