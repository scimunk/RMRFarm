package main

import (
	. "github.com/epixerion/RMRFarm/linker"
	. "github.com/epixerion/RMRFarm/rmrfarm"
)

type masterHandler struct {
	linker *ServerLinker
	clientList map[int32]Client
	updateMaster bool
}

func newMasterHandler() *masterHandler {
	masterHandler := &masterHandler{}
	masterHandler.clientList = make(map[int32]Client)
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
			rmrfarm.projectManager.HandleSendFile(ReadPacketSendFile(packet.(LargePacket)))
		case PACKET_RENDERFRAME:
			rmrfarm.projectManager.renderFrame(ReadPacketRenderFrame(packet))
		}
	}
	for _, clientState := range mh.linker.GetLastClientState() {
		if clientState.IsConnected {
			mh.clientList[clientState.ClientInt.GetId()] = clientState.ClientInt
			mh.updateMaster = true
		}else{
			delete(mh.clientList, clientState.ClientInt.GetId())
		}
	}

	if mh.updateMaster{
		mh.SendSlaveInfo()
		mh.updateMaster = false
	}
}

func (mh *masterHandler) SendSlaveInfo(){
	packet := &PacketSlaveInfo{}
	packet.SlaveName = rmrfarm.conf.SlaveName
	packet.Availlable = rmrfarm.projectManager.projectHook == nil
	for _, project := range rmrfarm.projectManager.projectList{
		if project.state == STATE_READY{
			packet.ProjectReady = append(packet.ProjectReady, project.projectName)
		}
	}
	mh.linker.SendPacketToAll(packet)
}