package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

// type Config struct {
// 	App      `json:"app"`
// 	Server   `json:"server"`
// 	Database `json:"database"`
// }

// type App struct {
// 	Name string `json:"name"`
// }

// type Server struct {
// 	Port int `json:"port"`
// }

// type Database struct {
// 	Host     string `json:"host"`
// 	Port     int    `json:"port"`
// 	Name     string `json:"name"`
// 	UserName string `json:"username"`
// 	Password string `json:"password"`
// }

type Config struct {
	App      app
	Server   server
	Database database
}

type app struct {
	Name   string
	Author string
	Email  string
}

type server struct {
	Port int
}

type database struct {
	Host     string
	Port     int
	Name     string
	UserName string
	Password string
}

const cfgPathTOML string = "./config.toml"

var Cfg Config

func init() {
	toml.DecodeFile(cfgPathTOML, &Cfg)
	log.Println("Config Init Over!")
}
