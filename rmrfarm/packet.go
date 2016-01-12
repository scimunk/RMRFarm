package rmrfarm

import (
	"github.com/epixerion/RMRFarm/linker"
)

const (
	PACKET_SLAVEINFO = iota
	PACKET_NEWPROJECT
	PACKET_REQUESTFILE
	PACKET_SENDFILE
	PACKET_SLAVEREADY
	PACKET_RENDERFRAME
)

/*
	Base PAcket
*/
type PacketData struct {
	PacketId byte
	Client   linker.Client
}

func readBasePacket(p linker.Packet) PacketData {
	return PacketData{p.GetId(), p.GetClient()}
}

func (p *PacketData) GetClient() linker.Client {
	return p.Client
}

func (p *PacketData) GetId() byte {
	return p.PacketId
}

type LargePacketData struct {
	PacketData
	Filepath string
	Size     int32
}

func (p *LargePacketData) GetFilePath() string {
	return p.Filepath
}

/*
	Packet Slave Info
*/
type PacketSlaveInfo struct {
	PacketData
	SlaveName string
}

func ReadPacketSlaveInfo(p linker.Packet) *PacketSlaveInfo {
	reader := CreateBinaryReader(p.GetBytes())
	packet := &PacketSlaveInfo{}
	packet.SlaveName = reader.ReadUtfString()
	return packet
}

func (p *PacketSlaveInfo) GetBytes() []byte {
	writer := CreateBinaryWriter()
	writer.WriteUtfString(p.SlaveName)
	return writer.Bytes()
}

/*
	Packet NewProject
*/
type PacketNewProject struct {
	LargePacketData
	ProjectName string
	FileData    []FileData
}

func ReadPacketNewProject(p linker.LargePacket) *PacketNewProject {
	reader := CreateBinaryReader(p.GetBytes())
	newpacket := &PacketNewProject{}
	newpacket.Filepath = p.GetFilePath()
	newpacket.PacketId = p.GetId()

	newpacket.ProjectName = reader.ReadUtfString()
	len := reader.ReadInt8()
	for i := int8(0); i < len; i++ {
		fileD := FileData{}
		fileD.File = reader.ReadUtfString()
		fileD.Path = reader.ReadUtfString()
		fileD.IsExterne = reader.ReadBool()
		newpacket.FileData = append(newpacket.FileData, fileD)
	}
	return newpacket
}

func (p *PacketNewProject) GetBytes() []byte {
	writer := CreateBinaryWriter()
	writer.WriteUtfString(p.ProjectName)
	writer.WriteInt8(byte(len(p.FileData)))
	for _, file := range p.FileData {
		writer.WriteUtfString(file.File)
		writer.WriteUtfString(file.Path)
		writer.WriteBool(file.IsExterne)
	}
	return writer.Bytes()
}

/*
	Packet Request File
*/
type PacketRequestFile struct {
	PacketData
	FileList []FileData
}

func ReadPacketRequestFile(packet linker.Packet) *PacketRequestFile {
	reader := CreateBinaryReader(packet.GetBytes())
	newpacket := &PacketRequestFile{PacketData: readBasePacket(packet)}
	len := reader.ReadInt8()
	for i := int8(0); i < len; i++ {
		fileD := FileData{}
		fileD.File = reader.ReadUtfString()
		fileD.Path = reader.ReadUtfString()
		fileD.IsExterne = reader.ReadBool()
		newpacket.FileList = append(newpacket.FileList, fileD)
	}
	return newpacket
}

func (p *PacketRequestFile) GetBytes() []byte {
	writer := CreateBinaryWriter()
	writer.WriteInt8(byte(len(p.FileList)))
	for _, file := range p.FileList {
		writer.WriteUtfString(file.File)
		writer.WriteUtfString(file.Path)
		writer.WriteBool(file.IsExterne)
	}
	return writer.Bytes()
}

/*
	Packet SendFile
*/
type PacketSendFile struct {
	LargePacketData
	Path     string
	FileName string
}

func ReadPacketSendFile(p linker.LargePacket) *PacketSendFile {
	reader := CreateBinaryReader(p.GetBytes())
	newpacket := &PacketSendFile{}
	newpacket.Filepath = p.GetFilePath()
	newpacket.PacketId = p.GetId()

	newpacket.Path = reader.ReadUtfString()
	newpacket.FileName = reader.ReadUtfString()
	return newpacket
}

func (p *PacketSendFile) GetBytes() []byte {
	writer := CreateBinaryWriter()
	writer.WriteUtfString(p.Path)
	writer.WriteUtfString(p.FileName)
	return writer.Bytes()
}


/*
	Packet Slave Ready
*/
type PacketSlaveReady struct {
	PacketData
	ProjectName []string
	Availlable bool
}

func ReadPacketSlaveReady(packet linker.Packet) *PacketSlaveReady {
	reader := CreateBinaryReader(packet.GetBytes())
	newpacket := &PacketSlaveReady{PacketData: readBasePacket(packet)}
	newpacket.Availlable = reader.ReadBool()
	newpacket.ProjectName = reader.ReadUtfStringArray()
	return newpacket
}

func (p *PacketSlaveReady) GetBytes() []byte {
	writer := CreateBinaryWriter()
	writer.WriteBool(p.Availlable)
	writer.WriteUtfStringArray(p.ProjectName)
	return writer.Bytes()
}

/*
	Packet RenderFrame
*/
type PacketRenderFrame struct {
	PacketData
	ProjectName string
	Camera string
	FrameId int32
}

func ReadPacketRenderFrame(packet linker.Packet) *PacketRenderFrame {
	reader := CreateBinaryReader(packet.GetBytes())
	newpacket := &PacketRenderFrame{PacketData: readBasePacket(packet)}
	newpacket.ProjectName = reader.ReadUtfString()
	newpacket.Camera = reader.ReadUtfString()
	newpacket.FrameId = reader.ReadInt32()
	return newpacket
}

func (p *PacketRenderFrame) GetBytes() []byte {
	writer := CreateBinaryWriter()
	writer.WriteUtfString(p.ProjectName)
	writer.WriteUtfString(p.Camera)
	writer.WriteInt32(p.FrameId)
	return writer.Bytes()
}