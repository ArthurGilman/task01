package main

import (
	"log"
	"sync"
	"task01/internal/api"

	"task01/internal/cache"
	"task01/internal/config"
	"task01/internal/db"

	"task01/internal/kafka"
	"task01/internal/logger"
)

func main() {
	err := logger.Init()
	if err != nil {
		log.Fatalf("Failed to create logger : %v", err)
	}

	config.Init()

	if err = cache.Init(); err != nil {
		logger.ErrorLog.Fatalf("%v", err)
	}

	if err = db.ConnectionPsql(); err != nil {
		logger.ErrorLog.Fatalf("%v", err)
		return
	}

	if err = db.InitStruct(); err != nil {
		logger.ErrorLog.Fatalf("%v", err)
		return
	}

	if err = kafka.Init(); err != nil {
		logger.ErrorLog.Fatalf("%v", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	defer wg.Done()
	defer wg.Done()

	go api.Init()
	go kafka.StartReading()

	wg.Wait()
}
