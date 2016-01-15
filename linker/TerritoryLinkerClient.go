package linker

import (
	"encoding/binary"
	"net"
	"time"
	"fmt"
	"io/ioutil"
	"io"
	"github.com/epixerion/RMRFarm/logger"
)

type LinkerData struct {
	isConnected bool
	address     string
	packetOut   chan Packet
	packetIn    chan Packet
	logger logger.Logger
}

func StartClientLinker(address string) *LinkerData {
	linker := &LinkerData{false, address, make(chan Packet, 10000), make(chan Packet, 10000), nil}
	go linker.linkerOperation()
	return linker
}

func (linkerdata *LinkerData) SetLogger(logger logger.Logger){
	linkerdata.logger = logger
}

func (linkerdata *LinkerData) GetPacket() (packetIn []Packet) {
	for i := 0; i < 10000; i++ {
		select {
		case packet := <-linkerdata.packetIn:
			packetIn = append(packetIn, packet)
		default:
			break
		}

	}
	return
}

func (LinkerData *LinkerData) IsConnected() bool{
	return LinkerData.isConnected
}

func (linkerdata *LinkerData) SendPacket(packet Packet) {
	linkerdata.packetOut <- packet
}

func (linker *LinkerData) linkerOperation() {

	for true {
		conn, err := net.Dial("tcp", linker.address)
		client := linkerClient{0,nil, conn}
		logMsg(linker.logger, logger.LOG_LOW, "LINKER", "starting linker client at address", linker.address)
		if err != nil {
			logWarn(linker.logger, logger.LOG_MEDIUM, "LINKER", "Error, Can't Connect to server, trying again soon", err)
			time.Sleep(time.Second * 15)
		} else {
			logMsg(linker.logger, logger.LOG_LOW, "LINKER", "Connected to server")
			defer conn.Close()
			linker.isConnected = true

			go readPacket(linker, conn)

			var packet Packet
			for true {
				for i := 0; i < 10000; i++ {
					select {
					case packet = <-linker.packetOut:
						client.SendPacket(packet)
					default:
						break
					}

				}
				if !linker.isConnected {
					logWarn(linker.logger, logger.LOG_MEDIUM, "LINKER", "Link Broken for ip", linker.address)
					break
				}
			}
		}
	}
}

func readPacket(linker *LinkerData, conn net.Conn) {
	var packet Packet
	for true {
		buffer := make([]byte, 2)
		conn.Read(buffer)
		fmt.Println("receiving packet of size", buffer)
		buffer = make([]byte, binary.BigEndian.Uint16(buffer))
		_, err := conn.Read(buffer)
		if err != nil {
			linker.isConnected = false
			fmt.Println("error with connection", err)
			return
		}
		fmt.Println(buffer)
		if buffer[0] >= 1 {
			size := binary.BigEndian.Uint32(buffer[1:5])
			logMsg(linker.logger, logger.LOG_LOW, "LINKER", "receiving large packet size :", size)
			fmt.Println("receiving packet of size : ", size)
			f, err := ioutil.TempFile("", "largePacket")
			if err != nil{
				fmt.Println("err copy:", err)
			}
			_, err = io.CopyN(f, conn,int64(size))
			if err != nil{
				logErr(linker.logger, logger.LOG_HIGH,"LINKER", "large packet error" ,err)
			}
			packet = &linkerLargePacketData{linkerPacketData{buffer[5], buffer[6:len(buffer)], nil},f.Name()}
			linker.packetIn <- packet
			f.Close()
			logMsg(linker.logger, logger.LOG_LOW, "LINKER", "large packet received")
		}else {
			packet = &linkerPacketData{buffer[1], buffer[2:len(buffer)],nil}
			linker.packetIn <- packet
		}
	}
}
