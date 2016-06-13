package main

import (
	"log"
	"time"
)

type Listener struct {
	logger *log.Logger
	Config *Config
	Pool   *SSDBPool
}

func NewListener(logger *log.Logger, config *Config, pool *SSDBPool) (listener *Listener, err error) {
	listener = &Listener{
		logger: logger,
		Config: config,
		Pool:   pool,
	}
	logger.Println(config)
	return
}

func (listener *Listener) Listen() {
	t := time.NewTimer(time.Duration(listener.Config.ReloadSecond) * time.Second)
	go func() {
		for _ = range t.C {
			listener.ReloadConfig()
			t.Reset(time.Duration(listener.Config.ReloadSecond) * time.Second)
		}
	}()
}

func (listener *Listener) ReloadConfig() {
	conf, err := LoadConfig()
	if err != nil {
		listener.logger.Println(err)
		return
	}
	listener.Config = conf
	listener.logger.Println("config reload success")
	listener.logger.Println(conf)
}
