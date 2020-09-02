package config

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Config ...
type Config struct {
	DBConn     *gorm.DB
	HttpClient *http.Client
}

// Env ...
type Env struct {
	HTTPPort                                   string
	DBHost, DBPort, DBUser, DBPassword, DBName string
	Username, Password, Organization, Domain   string

	Conn *gorm.DB
}

// GlobalEnv global environment
var GlobalEnv Env

// Init app config
func Init() *Config {
	cfg := Config{}
	var ok bool

	GlobalEnv.HTTPPort, ok = os.LookupEnv("HTTP_PORT")
	if !ok {
		panic("missing HTTP_PORT environment")
	}

	GlobalEnv.DBHost, ok = os.LookupEnv("DB_HOST")
	if !ok {
		panic("missing DB_HOST environment")
	}

	GlobalEnv.DBPort, ok = os.LookupEnv("DB_PORT")
	if !ok {
		panic("missing DB_PORT environment")
	}

	GlobalEnv.DBUser, ok = os.LookupEnv("DB_USER")
	if !ok {
		panic("missing DB_USER environment")
	}

	GlobalEnv.DBPassword, ok = os.LookupEnv("DB_PASSWORD")
	if !ok {
		panic("missing DB_PASSWORD environment")
	}

	GlobalEnv.DBName, ok = os.LookupEnv("DB_NAME")
	if !ok {
		panic("missing DB_NAME environment")
	}

	if cfg.DBConn == nil {
		db, err := gorm.Open("postgres", fmt.Sprintf("host=%s user=%s "+
			"password=%s port=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME")))
		if err != nil {
			panic("Database Connection Failed")
		}

		db.LogMode(true)

		GlobalEnv.Conn = db
		defer db.Close()
	}

	return &cfg
}

// Db ..
func Db() *gorm.DB {
	return GlobalEnv.Conn
}
