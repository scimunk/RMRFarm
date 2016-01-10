package main

import (
	//"os"
	"archive/zip"
	"bufio"
	"github.com/epixerion/RMRFarm/linker"
	. "github.com/epixerion/RMRFarm/rmrfarm"
	"io"
	"os"
	"path/filepath"
	"strings"
	"github.com/epixerion/RMRFarm/logger"
	"io/ioutil"
	"os/exec"
)

const (
	STATE_RENDERING = iota
	STATE_WAITINGFILE
	STATE_READY
)

type projectData struct {
	projectName string
	fileList    []FileData
	mainClient  linker.Client
	client      linker.Client
	state       int8
}

func newProject(data *PacketNewProject, client linker.Client) *projectData {
	mainLog.SetColor("\x1b[32m").LogMsg(logger.LOG_INFO, "PROJECT", "STARTING NEW PROJECT")
	projectData := &projectData{}
	projectData.projectName = data.ProjectName
	projectData.fileList = data.FileData
	projectData.mainClient = data.Client
	projectData.client = client
	projectData.DecompressProjectDataFile(data.Filepath)
	projectData.checkAndRequestRequiredFile()
	os.MkdirAll(filepath.Join(rmrfarm.conf.Workspace, "/renderman/", projectData.projectName, "images"), os.ModePerm)
	return projectData
}

func (pd *projectData) DecompressProjectDataFile(largefilepath string) {
	//err := os.MkdirAll(filepath.Join(rmrfarm.conf.Workspace, filepath, os.ModePerm)

	r, err := zip.OpenReader(largefilepath)
	if err != nil {

	}
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {

		}
		path := filepath.Join(rmrfarm.conf.Workspace, f.Name)
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
		ftoWrite, _ := os.Create(path)
		_, err = io.CopyN(ftoWrite, rc, int64(f.UncompressedSize64))
		if err != nil {

		}
		ftoWrite.Close()
		rc.Close()
	}
	defer r.Close()
}

func (pd *projectData) changeState(state int8) {
	pd.state = state
	//rmrfarm.masterHandler.linker.SendPacket(&packetSlaveInfo{})
}

func (pd *projectData) updateProject() {
	if pd.state == STATE_READY {
		pd.MakeCurrentPCCompatibleRIB()
		pd.state = STATE_RENDERING
		go pd.startRender()
	}
}

func (pd *projectData) startRender() {
	mainLog.SetColor("\x1b[1m\x1b[5m\x1b[31m").LogMsg(logger.LOG_INFO,"RENDER","STARTING RENDERING !")
	cmd := exec.Command("prman","-Progress",  "-cwd", rmrfarm.conf.Workspace, filepath.Join("renderman",pd.projectName,"/rib/0003/perspShape_Final.0003.rib"))
	stdout, err := cmd.StdoutPipe()
	stderr, err := cmd.StderrPipe()
	if err != nil {

	}
	if err := cmd.Start(); err != nil {

	}
	go reader(stdout)
	go reader(stderr)
	if err := cmd.Wait(); err != nil {

	}
}

func reader(io io.ReadCloser) {
	for {
		reader := bufio.NewReader(io)
		str, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		mainLog.SetColor("\x1b[34m").LogMsg(logger.LOG_INFO, "RENDER", string(str))
	}
}

func (pd *projectData) checkAndRequestRequiredFile() {
	mainLog.SetColor("\x1b[32m").LogMsg(logger.LOG_INFO, "PROJECT", "CHECKING REQUIRED FILES")
	fileData := pd.checkRequiredFile()
	if len(fileData) > 0 {
		for _, f := range fileData{
			mainLog.SetColor("\x1b[32m").LogMsg(logger.LOG_INFO, "PROJECT", filepath.Join(f.Path, f.File), "NOT EXISTING, ADDING TO REQUEST")
		}
		pd.client.GetConn().SendPacket(&PacketRequestFile{PacketData{PACKET_REQUESTFILE, pd.mainClient}, fileData})
	} else {
		mainLog.SetColor("\x1b[32m").LogMsg(logger.LOG_INFO, "PROJECT", "EVERY FILE HAVE BEEN RECEIVED")
		pd.changeState(STATE_READY)
	}
}

func (pd *projectData) checkRequiredFile() []FileData {
	var fileToRequest []FileData
	for _, file := range pd.fileList {
		if _, err := os.Stat(filepath.Join(rmrfarm.conf.Workspace, strings.Replace(file.Path, ":", "_", -1), file.File)); os.IsNotExist(err) {
			mainLog.SetColor("\x1b[32m").LogMsg(logger.LOG_INFO, "PROJECT", "Required file check :", file.File)
			fileToRequest = append(fileToRequest, file)
		}
	}
	return fileToRequest
}

func (pd *projectData) HandleSendFile(packet *PacketSendFile) {
	err := os.MkdirAll(filepath.Join(rmrfarm.conf.Workspace, strings.Replace(packet.Path, ":", "_", -1)), os.ModePerm)
	if err != nil {
		mainLog.SetColor("\x1b[31m").LogErr(logger.LOG_INFO, "PROJECT", "Could'nt not create directory :", err)
	}
	err = os.Rename(packet.Filepath, filepath.Join(rmrfarm.conf.Workspace, strings.Replace(packet.Path, ":", "_", -1), packet.FileName))
	if err != nil {
		mainLog.SetColor("\x1b[31m").LogErr(logger.LOG_INFO, "PROJECT", "Couldn't Move Temp File", err)
	}
	mainLog.SetColor("\x1b[32m").LogMsg(logger.LOG_INFO, "PROJECT", "FILE ",filepath.Join(rmrfarm.conf.Workspace, strings.Replace(packet.Path, ":", "_", -1), packet.FileName), "Have Been Received")
	fileData := pd.checkRequiredFile()
	if len(fileData) == 0 {
		mainLog.SetColor("\x1b[32m").LogMsg(logger.LOG_INFO, "PROJECT", "EVERY FILE HAVE BEEN RECEIVED")
		pd.changeState(STATE_READY)
	}
}

func (pd *projectData) MakeCurrentPCCompatibleRIB(){
	filepath.Walk(filepath.Join(rmrfarm.conf.Workspace, "renderman", pd.projectName), pd.replaceAbsPath)
}


func  (pd *projectData) replaceAbsPath(path string, info os.FileInfo, err error) error{
	filename := filepath.Base(info.Name())
	if !strings.Contains(filename, ".rib"){
		return nil
	}
	mainLog.SetColor("\x1b[32m").LogMsg(logger.LOG_INFO, "PROJECT", "Will convert rib for local path", path)

	input, _ := ioutil.ReadFile(path)
	lines := strings.Split(string(input), "\n")
	foundone := false
	for i, line := range lines {
		for _ , file := range pd.fileList{
			if file.IsExterne{
				if strings.Contains(line, file.Path) {
					replacewith := filepath.Join(rmrfarm.conf.Workspace, strings.Replace(file.Path, ":", "_", -1), file.File)
					lines[i] = strings.Replace(lines[i], file.Path, replacewith, -1)
					foundone = true
				}
			}
		}
	}
	if foundone{
		output := strings.Join(lines, "\n")
		err = ioutil.WriteFile(path, []byte(output), 0644)
		if err != nil {

		}
	}
	return nil
}
