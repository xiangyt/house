package config

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"sync"
)

type Config struct {
	Mode  string `json:"mode"`
	Mysql struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"mysql"`
}

var c *Config
var once sync.Once

func Get() *Config {
	once.Do(func() {
		c = LoadConfig()
	})
	return c
}

func LoadConfig() *Config {
	configPath := os.Getenv("CONF_PATH")
	if configPath == "" {
		configPath = "./config/config.json"
	}

	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		logrus.Fatalf("load config failed! path:%s, err:%s", configPath, err)
	}
	config := &Config{}
	err = json.Unmarshal(buf, config)
	if err != nil {
		logrus.Fatalf("decode config file failed! err:%s", err)
	}
	logrus.Info("load config success!")
	return config
}
