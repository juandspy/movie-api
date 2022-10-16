package main

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	wantConfig := Configuration{
		Host:     "testHost",
		Port:     1234,
		User:     "testUser",
		Password: "testPassword",
		DBName:   "testDBName",
		Driver:   "testDriver",
	}

	os.Setenv("MOVIE_API__HOST", "testHost")
	os.Setenv("MOVIE_API__PORT", "1234")
	os.Setenv("MOVIE_API__USER", "testUser")
	os.Setenv("MOVIE_API__PASSWORD", "testPassword")
	os.Setenv("MOVIE_API__DBNAME", "testDBName")
	os.Setenv("MOVIE_API__DRIVER", "testDriver")

	gotConfig := LoadConfig()

	if gotConfig != wantConfig {
		t.Errorf("got %v want %v", gotConfig, wantConfig)
	}

	os.Setenv("MOVIE_API__PORT", "not an int")
	gotConfig = LoadConfig()
	if gotConfig.Port != 5432 {
		t.Errorf("an invalid port shouldn't update the default value. Got %v instead of %v", gotConfig.Port, 5432)
	}
}
