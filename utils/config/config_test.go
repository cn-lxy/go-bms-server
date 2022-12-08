package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestConfigOfJSON(t *testing.T) {
	var cfg Config
	file, err := os.ReadFile("../../config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("App Name: %v\n", cfg.App.Name)
	fmt.Printf("Server Port: %v\n", cfg.Server.Port)
	fmt.Printf("Database Host: %v\n", cfg.Database.Host)
	fmt.Printf("Database Port: %v\n", cfg.Database.Port)
	fmt.Printf("Database Name: %v\n", cfg.Database.Name)
	fmt.Printf("Database Username: %v\n", cfg.Database.UserName)
	fmt.Printf("Database Password: %v\n", cfg.Database.Password)
}

type db struct {
	Name string
	Host string
	Port int
}

type config struct {
	Db db
}

func TestConfigOfToml(t *testing.T) {
	var cfg config
	file, err := os.ReadFile("./test.toml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(file))
	toml.Decode(string(file), &cfg)
	fmt.Printf("cfg: %#v\n", cfg)
}
