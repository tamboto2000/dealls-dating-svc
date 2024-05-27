package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const defCfgPath = "./configs/app.yaml"

type Database struct {
	Driver   string `yaml:"driver"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
}

type SMTP struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	SenderName string `yaml:"sender_name"`
	Email      string `yaml:"email"`
	Password   string `yaml:"password"`
}

type Password struct {
	MinLength int    `yaml:"min_length"`
	MaxLength int    `yaml:"max_length"`
	Salt      string `yaml:"salt"`
}

type JWT struct {
	SignMethod string `yaml:"sign_method"`
	Key        string `yaml:"key"`
	Expiration int64  `yaml:"expiration"`
}

func (j JWT) ExpireInSeconds() time.Duration {
	return time.Second * time.Duration(j.Expiration)
}

type Auth struct {
	Password Password `yaml:"password"`
	JWT      JWT      `yaml:"jwt"`
}

type Cache struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func (c Cache) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

type PremiumFeatures struct {
	UnlimitedSwipe int64 `yaml:"unlimited_swipe"`
}

type Config struct {
	Database        Database        `yaml:"database"`
	SMTP            SMTP            `yaml:"smtp"`
	Auth            Auth            `yaml:"auth"`
	Cache           Cache           `yaml:"cache"`
	PremiumFeatures PremiumFeatures `yaml:"premium_features"`
}

func LoadConfig() (Config, error) {
	var cfg Config

	f, err := os.Open(defCfgPath)
	if err != nil {
		return cfg, err
	}

	defer f.Close()

	err = yaml.NewDecoder(f).Decode(&cfg)

	return cfg, err
}
