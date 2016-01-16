package linker

import (
	"net"
	"fmt"
	"io"
	"os"
)

type Client interface{
	GetId() int32
	sendPacket(Packet)
}

type linkerClient struct {
	clientId   int32
	servRef *ServerLinker
	clientConn net.Conn
}

func (linkerClient *linkerClient) GetId() int32{
	return linkerClient.clientId
}

func (client *linkerClient) sendPacket(packet Packet) {
	largePacket, isLargePacket := packet.(LargePacket)
	headersize := 4
	if isLargePacket{
		headersize = 8
	}
	length := uint16(len(packet.GetBytes()) + headersize - 2)
	data := make([]byte, headersize)
	data[0] = uint8(length >> 8)
	data[1] = uint8(length & 0xff)
	data[2] = 0
	if isLargePacket {
		f, err := os.Stat(largePacket.GetFilePath())
		if err != nil{
			fmt.Println("large packet err ",err)
		}
		data[2] = 1
		data[3] = uint8(f.Size() >> 24)
		data[4] = uint8(f.Size() >> 16)
		data[5] = uint8(f.Size() >> 8)
		data[6] = uint8(f.Size()  & 0xff)
	}
	data[headersize - 1] = byte(packet.GetId())
	data = append(data, packet.GetBytes()...)
	client.clientConn.Write(data)
	if isLargePacket{
		f, err := os.Open(largePacket.GetFilePath())
		if err != nil {
			fmt.Println("open error", err)
		}
		size, _ := f.Stat()
		_, err = io.CopyN(client.clientConn, f, size.Size())
		if err != nil {
			fmt.Println("open error", err)
		}
		f.Close()
		client.clientConn.Write([]byte{})
	}
}