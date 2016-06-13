package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	ReloadSecond int    `json:"reloadSecond"`
	Args         Args   `json:"args"`
	Master       Node   `json:"master"`
	Nodes        []Node `json:"nodes"`
}

type Node struct {
	Id       string `json:"id"`
	Route    string `json:"route"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Open     bool   `json:"open"`
	OwnArgs  bool   `json:"ownArgs"`
	Args     Args   `json:"args"`
}

type Args struct {
	GetClientTimeout int `json:"getClientTimeout"`
	MaxPoolSize      int `json:"maxPoolSize"`
	MinPoolSize      int `json:"minPoolSize"`
	AcquireIncrement int `json:"acquireIncrement"`
	MaxIdleTime      int `json:"maxIdleTime"`
	MaxWaitSize      int `json:"maxWaitSize"`
	HealthSecond     int `json:"healthSecond"`
}

func LoadConfig() (config *Config, err error) {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		return
	}
	config = new(Config)
	err = json.Unmarshal(file, config)
	return
}
