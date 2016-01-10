package main

import (
	"fmt"
	. "github.com/epixerion/RMRFarm/rmrfarm"
	. "github.com/epixerion/RMRFarm/linker"
	"github.com/epixerion/RMRFarm/logger"
	"path/filepath"
)

type slaveData struct {
	id             int8
	linker         *LinkerData
	slaveName      string
	available      bool
	project        *projectData
	workingOnFrame int32
}

func StartSlave(ip string) *slaveData {
	slaveData := &slaveData{}
	slaveData.linker = StartClientLinker(ip)
	slaveData.linker.SetLogger(mainLog)
	return slaveData
}

func (slave *slaveData) UpdateSlave() {
	for _, packet := range slave.linker.GetPacket() {
		switch packet.GetId() {
		case PACKET_SLAVEINFO:
			fmt.Println("received slave info")
			slave.readSlaveData(ReadPacketSlaveInfo(packet))
		case PACKET_REQUESTFILE:
			slave.handleSlaveFileRequest(ReadPacketRequestFile(packet))
		}
	}
}

func (slave *slaveData) AssignSlaveToProject(project *projectData) {
	slave.project = project
	packet := &PacketNewProject{}
	packet.PacketId = PACKET_NEWPROJECT
	packet.Filepath = filepath.Join(rmrfarm.conf.MayaWorkspace, "RMRFarm", project.projectName+".zip")

	packet.FileData = project.fileData
	packet.ProjectName = project.projectName
	packet.Camera = "nil"
	slave.SendPacket(packet)
}

func (slave *slaveData) readSlaveData(packet *PacketSlaveInfo) {
	slave.slaveName = packet.SlaveName
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
		filepacket.FileName = file.File
		filepacket.Path = file.Path
		slave.linker.SendPacket(filepacket)
	}
}

func (slave *slaveData) SendPacket(packet Packet) {
	slave.linker.SendPacket(packet)
}

func (slave *slaveData) SendFile(filepath string) {
	/*folders, _ := ioutil.ReadDir(filepath.Join('./renderman/', NOMPROJECT))

	regex, _ := regexp.Compile("/\"(.*?)\"/g")
	workingdir, _ := os.Getwd()

	for _, dir := range folders {
		if dir.IsDir() && string.hasPrefix(dir.Name(), 'frame') {
		files, _ := ioutil.ReadDir(filepath.Join('./renderman/', NOMPROJECT, dir.Name()))

		for _, doc := range files {
		if !strings.HasSuffix(doc.Name(), ".rib")
		continue
		content, _ := ioutil.ReadFile(filepath.Join('./renderman/', NOMPROJECT, dir.Name(), doc.Name()))


		for quoted := regex.FindAllString(string(content[:]), -1) {
		if strings.Index(quoted, "/") != -1 && strings.Index(quoted, "@") == -1 && strings.Index(quoted, ":") == -1  {
		if !strings.HasSuffix(quoted, "/") {
		if path.IsAbs(quoted) {
		rel := filepath.Rel(workingdir, quoted)
		} else {
		abs := path.Abs(quoted)
		}
		}
		}
		}
		}
		}
	}*/
}
