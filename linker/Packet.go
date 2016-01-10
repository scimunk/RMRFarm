package linker

import (

)

type Packet interface{
	GetId() byte
	GetBytes() []byte
	GetClient() Client
}

type LargePacket interface{
	Packet
	GetFilePath() string
}

type linkerPacketData struct {
	Packetid byte
	Bytedata []byte
	Client   *linkerClient
}

func (packetD *linkerPacketData) GetId() byte{
	return packetD.Packetid
}

func (packetD *linkerPacketData) GetBytes() []byte{
	return packetD.Bytedata
}

func (packetD *linkerPacketData) GetClient() Client{
	return packetD.Client
}

type linkerLargePacketData struct {
	linkerPacketData
	filepath string
}

func (packetD *linkerLargePacketData) GetFilePath() string{
	return packetD.filepath
}