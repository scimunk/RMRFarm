package main

import (
	"time"
	"github.com/epixerion/RMRFarm/logger"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var rmrfarm RMRFarm
var mainLog logger.Logger

type RMRFarm struct {
	conf           *config
	masterHandler  *masterHandler
	projectManager *projectManager
}

func main() {
	mainLog = logger.NewLogger(logger.LOG_HIGH)
	rmrfarm.conf = loadConfiguration()
	rmrfarm.masterHandler = newMasterHandler()
	rmrfarm.masterHandler.updateMasterHandler()
	rmrfarm.projectManager = startProjectManager()
	for {
		rmrfarm.masterHandler.updateMasterHandler()
		rmrfarm.projectManager.updateProjectManager()
		time.Sleep(time.Second)
	}
}
