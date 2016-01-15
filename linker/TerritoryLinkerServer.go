package linker

import (
	"encoding/binary"
	"fmt"
	"net"
	"io"
	"io/ioutil"
	"github.com/epixerion/RMRFarm/logger"
)

/*
	The QuantumLinkerServer  Link the Connection/Game Server
*/
type ServerLinker struct {
	address       string
	clientIdPool  int32
	packetOut     chan Packet
	packetOutAll     chan Packet
	packetIn      chan Packet
	logger logger.Logger
	clienthandler chan LinkerClientHandler
}

type LinkerClientHandler struct {
	IsConnected bool
	ClientInt   Client
}

func check(err error){
	if err!=nil{
		panic(err)
	}
}

func (linker *ServerLinker) GetPacket() (packetIn []Packet) {
	for i := 0; i < 10000; i++ {
		select {
		case packet := <-linker.packetIn:
			packetIn = append(packetIn, packet)
		default:
			break
		}

	}
	return
}

func (linkerdata *ServerLinker) SetLogger(logger logger.Logger){
	linkerdata.logger = logger
}

func (linker *ServerLinker) GetLastClientState() (lastClient []LinkerClientHandler) {
	for i := 0; i < 1000; i++ {
		select {
		case client := <-linker.clienthandler:
			lastClient = append(lastClient, client)
		default:
			break
		}

	}
	return
}

//function used to send a packet to the client specified in packet data
func (linker *ServerLinker) SendPacket(packet Packet) {
	linker.packetOut <- packet
}

//function used to send a packet to the client specified in packet data
func (linker *ServerLinker) SendPacketToAll(packet Packet) {
	linker.packetOut <- packet
}

func StartServerLinker(address string) *ServerLinker {
	linker := &ServerLinker{address, 0, make(chan Packet, 100), make(chan Packet, 100), make(chan Packet, 100),nil, make(chan LinkerClientHandler, 100)}
	go linker.HandleLinker()
	return linker
}

func (linker *ServerLinker) HandleLinker() {
	fmt.Println("starting linker handler server at address", linker.address)

	clientchannel := make(chan LinkerClientHandler)

	clientListener, err := net.Listen("tcp", linker.address)
	check(err)

	defer clientListener.Close()

	go linker.connectionHandler(clientchannel, clientListener)

	var newclient LinkerClientHandler
	var packet Packet
	for true {

		for i := 0; i < 10; i++ {

			select {
			case newclient = <-clientchannel:
				linker.clienthandler <- newclient
			default:
				break
			}
		}

		for i := 0; i < 50; i++ {
			select {
			case packet = <-linker.packetOut:
				packet.GetClient().GetConn().SendPacket(packet)
			default:
				break
			}
		}
	}
}

func (linker *ServerLinker) connectionHandler(clientChannel chan LinkerClientHandler, ln net.Listener) {
	for true {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		client := linkerClient{linker.getNewId(), linker, conn}
		go client.handleConnection(clientChannel, linker.packetIn)
		fmt.Println("Receiving new Linker connection from port ", linker.address)
	}
}

func (linker *ServerLinker) getNewId() int32 {
	linker.clientIdPool++
	return linker.clientIdPool
}

func (client *linkerClient) handleConnection(clientchannel chan LinkerClientHandler, packetIn chan Packet) {
	defer client.handleDeconnection(clientchannel)

	//waiting for login information
	var packet Packet
	clientchannel <- LinkerClientHandler{true, client}
	//while for packet reception
	for true {
		buffer := make([]byte, 2)
		client.clientConn.Read(buffer)
		buffer = make([]byte, binary.BigEndian.Uint16(buffer))
		_, err := client.clientConn.Read(buffer)
		if err != nil {
			return
		}
		if buffer[0] >= 1 {
			size := binary.BigEndian.Uint32(buffer[1:5])
			logMsg(client.servRef.logger, logger.LOG_LOW, "LINKER", "receiving large packet size :", size)
			f, err := ioutil.TempFile("", "largePacket")
			if err != nil{
				fmt.Println("err copy:", err)
			}

			_, err = io.CopyN(f, client.clientConn,int64(size))
			if err != nil{
				fmt.Println("err copy:", err)
			}
			packet = &linkerLargePacketData{linkerPacketData{buffer[5], buffer[6:len(buffer)], client},f.Name()}
			packetIn <- packet
			f.Close()
			logMsg(client.servRef.logger, logger.LOG_LOW, "LINKER", "Large Packet SuccessFully Received", f.Name())
		}else {
			packet = &linkerPacketData{buffer[1], buffer[2:len(buffer)], client}
			packetIn <- packet
		}
	}
}

func (client *linkerClient) handleDeconnection(clientchannel chan LinkerClientHandler) {
	clientchannel <- LinkerClientHandler{false, client}
	client.clientConn.Close()
}
