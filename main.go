package main

import (
	"log"
	"os"
)

func main() {
	logFile, err := os.OpenFile("project.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	logger := log.New(logFile, "", log.LstdFlags)

	config, err := LoadConfig()
	if err != nil {
		logger.Fatalln(err)
		return
	}
	logger.Println("config load success")
	logger.Println(config)

	pool := NewSSDBPool(logger, config)
	pool.Start()

	server := NewServer(logger, config, pool)
	server.Start()
}
