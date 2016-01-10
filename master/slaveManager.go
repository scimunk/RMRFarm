package main

import (
	"fmt"
	. "github.com/epixerion/RMRFarm/linker"
)

type slaveManager struct {
	slaveData []*slaveData
}

func newSlaveManager() *slaveManager {
	slaveManager := &slaveManager{}
	go slaveManager.StartSlaveManager()
	return slaveManager
}

func (slaveM *slaveManager) StartSlaveManager() {
	for id, ip := range rmrfarm.conf.SlaveListIp {
		slave := StartSlave(ip)
		slave.id = int8(id)
		slaveM.slaveData = append(slaveM.slaveData, slave)
	}
}

func (slaveM *slaveManager) updateSlaveManager() {
	for _, slave := range slaveM.slaveData {
		slave.UpdateSlave()
	}
}

func (slaveM *slaveManager) getAvaillableSlave() []*slaveData {
	return slaveM.slaveData
}

func (slaveM *slaveManager) SendPacketToAll(packet Packet) {
	for _, slave := range slaveM.slaveData {
		fmt.Println(slave.id)
		slave.SendPacket(packet)
	}
}
