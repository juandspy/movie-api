package main

import (
	"log"
	"os"
	"strconv"
)

type Configuration struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	Driver   string
}

// loadConfig replace the default configuration with values from env variables
func LoadConfig() Configuration {
	// default configuration
	config := Configuration{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "mysecretpassword",
		DBName:   "postgres",
		Driver:   "postgres",
	}

	// replace with env variables
	host := os.Getenv("MOVIE_API__HOST")
	if host != "" {
		config.Host = host
	}
	port := os.Getenv("MOVIE_API__PORT")
	if port != "" {
		intPort, err := strconv.Atoi(port)
		if err != nil {
			log.Println("invalid value at 'MOVIE_API__PORT'. Using default port:", config.Port)
		} else {
			config.Port = intPort
		}
	}
	user := os.Getenv("MOVIE_API__USER")
	if user != "" {
		config.User = user
	}
	password := os.Getenv("MOVIE_API__PASSWORD")
	if password != "" {
		config.Password = password
	}
	dbName := os.Getenv("MOVIE_API__DBNAME")
	if dbName != "" {
		config.DBName = dbName
	}
	driver := os.Getenv("MOVIE_API__DRIVER")
	if driver != "" {
		config.Driver = driver
	}

	return config
}
