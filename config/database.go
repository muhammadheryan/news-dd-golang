package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
)

var config = ConfigDB{}

// ConfigDB db seting
type ConfigDB struct {
	User          string
	Password      string
	Host          string
	Port          string
	Dbname        string
	Redisaddr     string
	Redispassword string
	Redisdb       int
}

// ConnectDB returns initialized gorm.DB
func ConnectDB() (*gorm.DB, error) {
	config.Read()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.User, config.Password, config.Host, config.Port, config.Dbname)

	db, err := gorm.Open("mysql", dsn)
	db.SingularTable(true)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// ConnectDB returns initialized gorm.DB
func ConnectRedis() *redis.Client {
	config.Read()

	client := redis.NewClient(&redis.Options{
		Addr:     config.Redisaddr,
		Password: config.Redispassword,
		DB:       config.Redisdb,
	})

	return client
}

// Read and parse the configuration file
func (c *ConfigDB) Read() {
	if _, err := toml.DecodeFile("config.toml", &c); err != nil {
		log.Fatal(err)
	}
}
