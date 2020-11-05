package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config 配置
type Config struct {
	File   string
	Log    int `yaml:"log"`
	Server struct {
		ID        int64  `yaml:"id"`
		Host      string `yaml:"host"`
		Port      int    `yaml:"port"`
		JWTSecret string `yaml:"secret"`
		Expire    int64  `yaml:"expire"`
	}
	Mysql struct {
		Host         string `yaml:"host"`
		Port         string `yaml:"port"`
		User         string `yaml:"user"`
		Password     string `yaml:"password"`
		Database     string `yaml:"database"`
		MaxIdleConns int    `yaml:"max_idle_conns"`
		MaxOpenConns int    `yaml:"max_open_conns"`
		MaxLifeTime  int64  `yaml:"max_life_time"`
		ShowSQL      bool   `yaml:"show_sql"`
		Sync         bool   `yaml:"sync"`
		UseCache     bool   `yaml:"use_cache"`
	}
	Cache  string `yaml:"cache"`
	Memory struct {
		Size int `yaml:"size"`
	}
	Redis struct {
		Host     string `yaml:"host"`
		Password string `yaml:"password"`
		Port     int    `yaml:"port"`
		DB       int    `yaml:"db"`
	}
}

// Load 加载
func Load(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	config.File = file
	return &config, nil
}
