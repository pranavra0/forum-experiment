package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App struct {
		Port string `yaml:"addr"`
		Env  string `yaml:"env"`
	} `yaml:"app"`

	DB struct {
		Path          string `yaml:"path"`
		MaxOpenConns  int    `yaml:"max_open_conns"`
		BusyTimeoutMS int    `yaml:"busy_timeout_ms"`
		ForeignKeys   bool   `yaml:"foreign_keys"`
	} `yaml:"db"`

	Templates struct {
		Dir  string `yaml:"dir"`
		Base string `yaml:"base"`
	} `yaml:"templates"`

	Bootstrap struct {
		Admins []struct {
			Username string `yaml:"username"`
			Email    string `yaml:"email"`
		} `yaml:"admins"`
	} `yaml:"bootstrap"`

	Pagination struct {
		ThreadsPerPage int `yaml:"threads_per_page"`
		RepliesPerPage int `yaml:"replies_per_page"`
	} `yaml:"pagination"`

	Features struct {
		AllowRegistration   bool `yaml:"allow_registration"`
		AllowAnonymous      bool `yaml:"allow_anonymous"`
		EnableAPI           bool `yaml:"enable_api"`
		AllowThreadDeletion bool `yaml:"allow_thread_deletion"`
	} `yaml:"features"`
}

var C Config

func LoadConfig() {
	_ = godotenv.Load()

	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("❌ Error reading config.yaml: %v", err)
	}

	if err := yaml.Unmarshal(data, &C); err != nil {
		log.Fatalf("❌ Error parsing config.yaml: %v", err)
	}

	log.Println("✅ Loaded configuration")

	if envDB := os.Getenv("DB_PATH"); envDB != "" {
		C.DB.Path = envDB
	}
}
