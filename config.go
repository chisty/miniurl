package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type AWS struct {
	Table     string `json:"aws_table"`
	Region    string `json:"aws_region"`
	AccessKey string
	SecretKey string
}

type Redis struct {
	Host      string `json:"host"`
	TTL       int    `json:"ttl"`
	MaxIdle   int    `json:"max_idle"`
	MaxActive int    `json:"max_active"`
}

type Config struct {
	AWS         AWS   `json:"aws"`
	Redis       Redis `json:"redis"`
	Port        int   `json:"port"`
	IdGenNodeId int   `json:"idgen_node_id"`
}

func defaultConfig() Config {
	return Config{
		AWS: AWS{
			Table:  "eatigo1",
			Region: "ap-southeast-1",
		},
		Redis: Redis{
			Host:      "localhost:6379",
			TTL:       10,
			MaxIdle:   80,
			MaxActive: 12000,
		},
		Port:        9000,
		IdGenNodeId: 1,
	}
}

func LoadConfig() Config {
	// accessKey := os.Getenv("AWS_ACCESS_KEY")
	// if len(accessKey) == 0 {
	// 	panic("aws access key cannot be empty.")
	// }

	// secretKey := os.Getenv("AWS_SECRET_KEY")
	// if len(secretKey) == 0 {
	// 	panic("aws access key cannot be empty.")
	// }

	var c Config
	f, err := os.Open(".config")
	if err != nil {
		c = defaultConfig()
	} else {
		dec := json.NewDecoder(f)
		err = dec.Decode(&c)
		if err != nil {
			panic(err)
		}
	}

	c.AWS.AccessKey = os.Getenv("AWS_ACCESS_KEY")
	if len(c.AWS.AccessKey) == 0 {
		panic("aws access key cannot be empty.")
	}

	c.AWS.SecretKey = os.Getenv("AWS_SECRET_KEY")
	if len(c.AWS.SecretKey) == 0 {
		panic("aws access key cannot be empty.")
	}

	table := os.Getenv("AWS_TABLE")
	if len(table) > 0 {
		c.AWS.Table = table
	}
	if len(c.AWS.Table) == 0 {
		panic("aws dynamo db table name cannot be empty.")
	}

	region := os.Getenv("AWS_REGION")
	if len(region) > 0 {
		c.AWS.Region = region
	}
	if len(c.AWS.Region) == 0 {
		panic("aws region name cannot be empty.")
	}

	host := os.Getenv("Redis_URL")
	if len(host) > 0 {
		c.Redis.Host = host
	}
	if len(c.Redis.Host) == 0 {
		panic("redis url cannot be empty.")
	}

	ttl := os.Getenv("Redis_TTL")
	if len(ttl) > 0 {
		c.Redis.TTL, err = strconv.Atoi(ttl)
		if err != nil {
			panic("redis ttl is invalid.")
		}
	}

	if c.Redis.TTL <= 0 {
		panic("redis ttl is invalid.")
	}

	fmt.Println("Successfully loaded configuration.")
	return c
}
