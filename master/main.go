package main

import (
	"github.com/epixerion/RMRFarm/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var rmrfarm RMRFarm
var mainLog logger.Logger
var exit = true

type RMRFarm struct {
	conf           *config
	slaveManager   *slaveManager
	projectManager *projectManager
}

func exitHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		exit = false
	}()
}

func main() {
	mainLog = logger.NewLogger(logger.LOG_HIGH)
	rmrfarm.conf = loadConfiguration()
	rmrfarm.slaveManager = newSlaveManager()
	rmrfarm.projectManager = startProjectManager()
	go handleCommand()
	go exitHandler()
	defer rmrfarm.projectManager.saveProject()
	for exit {
		rmrfarm.slaveManager.updateSlaveManager()
		rmrfarm.projectManager.updateProjectManager()
		time.Sleep(time.Second * 1)
		//time.Sleep(time.Second * )
	}
}
