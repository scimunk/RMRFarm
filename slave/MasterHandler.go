package main

import (
	"fmt"
	. "github.com/epixerion/RMRFarm/linker"
	. "github.com/epixerion/RMRFarm/rmrfarm"
)

type masterHandler struct {
	linker *ServerLinker
}

func newMasterHandler() *masterHandler {
	masterHandler := &masterHandler{}
	masterHandler.linker = StartServerLinker(rmrfarm.conf.Ip)
	masterHandler.linker.SetLogger(mainLog)
	return masterHandler
}

func (mh *masterHandler) updateMasterHandler() {
	for _, packet := range mh.linker.GetPacket() {
		switch packet.GetId() {
		case PACKET_NEWPROJECT:
			rmrfarm.projectManager.startProject(ReadPacketNewProject(packet.(LargePacket)), packet.GetClient())
		case PACKET_SENDFILE:
			rmrfarm.projectManager.currentProject.HandleSendFile(ReadPacketSendFile(packet.(LargePacket)))
		case PACKET_RENDERFRAME:
			rmrfarm.projectManager.renderFrame(ReadPacketRenderFrame(packet))
		}
	}
	for _, clientState := range mh.linker.GetLastClientState() {
		fmt.Println("sending slave info")
		clientState.ClientInt.GetConn().SendPacket(&PacketSlaveInfo{PacketData{PACKET_SLAVEINFO, clientState.ClientInt}, "truc"})
	}
}
