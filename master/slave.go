package main

import (
	. "github.com/epixerion/RMRFarm/rmrfarm"
	. "github.com/epixerion/RMRFarm/linker"
	"github.com/epixerion/RMRFarm/logger"
	"path/filepath"
)

type slaveData struct {
	id             int8
	ip string
	linker         *LinkerData
	slaveName      string
	available      bool
	project        *projectData
	workingOnFrame int32
	projectReady []string
}

func StartSlave(ip string) *slaveData {
	slaveData := &slaveData{}
	slaveData.linker = StartClientLinker(ip)
	slaveData.ip = ip
	slaveData.linker.SetLogger(mainLog)
	return slaveData
}

func (slave *slaveData) UpdateSlave() {
	for _, packet := range slave.linker.GetPacket() {
		switch packet.GetId() {
		case PACKET_SLAVEINFO:
			slave.readSlaveData(ReadPacketSlaveInfo(packet))
		case PACKET_REQUESTFILE:
			slave.handleSlaveFileRequest(ReadPacketRequestFile(packet))
		case PACKET_FRAMECOMPLETED:
			slave.handleFrameCompleted( ReadPacketFrameCompleted(packet.(LargePacket)))
		}
	}
}

func (slave *slaveData) AssignSlaveToProject(project *projectData) {
	slave.project = project
	packet := &PacketNewProject{}
	packet.PacketId = PACKET_NEWPROJECT
	packet.Filepath = filepath.Join(rmrfarm.conf.MayaWorkspace, "RMRFarm", project.ProjectName +".zip")
	packet.FileData = project.FileData
	packet.ProjectName = project.ProjectName
	slave.SendPacket(packet)
}

func (slave *slaveData) readSlaveData(packet *PacketSlaveInfo) {
	mainLog.LogMsg(logger.LOG_LOW, "PROJECT", "SLAVE UPDATE INFO", packet.SlaveName)
	slave.slaveName = packet.SlaveName
	slave.projectReady = packet.ProjectReady
	slave.available = packet.Availlable
}

func (slave *slaveData) handleSlaveFileRequest(packet *PacketRequestFile) {
	for _, file := range packet.FileList {
		filepacket := &PacketSendFile{}
		filepacket.PacketId = PACKET_SENDFILE
		mainLog.LogMsg(logger.LOG_INFO, "PROJECT", "received file request")
		if !file.IsExterne {
			filepacket.Filepath = filepath.Join(rmrfarm.conf.MayaWorkspace, file.Path, file.File)
		} else {

			filepacket.Filepath = filepath.Join(file.Path, file.File)
		}
		filepacket.IsExterne = file.IsExterne
		filepacket.FileName = file.File
		filepacket.Path = file.Path
		slave.linker.SendPacket(filepacket)
	}
}

func (slave *slaveData) SendPacket(packet Packet) {
	slave.linker.SendPacket(packet)
}

func (slave *slaveData) StartRenderFrame(project *projectData, frame *frame){
	packet := &PacketRenderFrame{}
	packet.ProjectName = project.ProjectName
	packet.Camera = project.Camera[project.cameraToRender]
	packet.FrameId = frame.frameId
	frame.renderedBy = slave.slaveName
	frame.state = FRAMESTATE_RENDERING
	slave.available = false
	slave.linker.SendPacket(packet)
}

func (slave *slaveData) handleFrameCompleted(packet *PacketFrameCompleted){
	slave.project.handleFrameCompleted(packet)
}

func (slave *slaveData) printSlaveInfo(){
	str := ""
	if slave.slaveName == ""{
		str += slave.ip
	}else{
		str += slave.slaveName
	}
	if (slave.linker.IsConnected()) {
		mainLog.SetColor(logger.COLOR_GREEN).LogMsg(logger.LOG_INFO, "SLAVE", str + " Is Connected")
	}
	if (!slave.linker.IsConnected()) {
		mainLog.SetColor(logger.COLOR_RED).LogMsg(logger.LOG_INFO, "SLAVE", str + " Is Not Connected")
	}
}