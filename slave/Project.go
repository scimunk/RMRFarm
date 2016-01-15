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
	"strconv"
)

const (
	STATE_RENDERING = iota
	STATE_WAITINGFILE
	STATE_READY
)

type projectData struct {
	projectName string
	fileList    []FileData
	client      linker.Client
	state       int8
}

func newProject(data *PacketNewProject, client linker.Client) *projectData {
	mainLog.SetColor(logger.COLOR_GREEN).LogMsg(logger.LOG_INFO, "PROJECT", "STARTING NEW PROJECT")
	projectData := &projectData{}
	projectData.projectName = data.ProjectName
	projectData.fileList = data.FileData
	projectData.client = client
	projectData.state = STATE_WAITINGFILE
	projectData.DecompressProjectDataFile(data.Filepath)
	projectData.CompatibiliseRIB()
	return projectData
}

func (pd *projectData) DecompressProjectDataFile(largefilepath string) {
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
	if pd.state == STATE_READY {
		rmrfarm.masterHandler.SendSlaveInfo()
	}
}

func (pd *projectData) updateProject() {
	if pd.state == STATE_WAITINGFILE{
		pd.checkAndRequestRequiredFile()
	}
}

func (pd *projectData) startRender(camera string, frameId int32) {
	os.MkdirAll(filepath.Join(rmrfarm.conf.Workspace, "/renderman/", pd.projectName, "images"), os.ModePerm)
	mainLog.SetColor(logger.COLOR_RED).LogMsg(logger.LOG_INFO,"RENDER","STARTING RENDERING !")
	cmd := exec.Command("prman","-Progress",  "-cwd", rmrfarm.conf.Workspace,
		filepath.Join("renderman",pd.projectName,"/rib/000" + strconv.Itoa(int(frameId))+ "/"+camera + "Shape_Final.000"+strconv.Itoa(int(frameId))+".rib"))
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "RMSPROJ_FROM_ENV="+rmrfarm.conf.Workspace)
	cmd.Env = append(cmd.Env, "RMSPROJ="+rmrfarm.conf.Workspace)
	stdout, err := cmd.StdoutPipe()
	stderr, err := cmd.StderrPipe()
	if err != nil {
			//prman -cwd C:\Users\epixe\Desktop\GoServer renderman/Andie_Furv2/rib/job/0003/perspShape_Final.0003.rib
	}
	if err := cmd.Start(); err != nil {

	}
	go pd.frameReader(stderr,frameId)
	go reader(stdout)
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
		mainLog.SetColor(logger.COLOR_BLUE).LogMsg(logger.LOG_INFO, "RENDER", string(str))
	}
}

func (pd *projectData) frameReader(io io.ReadCloser, frameid int32){
	completed := false
	for {
		reader := bufio.NewReader(io)
		str, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		mainLog.SetColor(logger.COLOR_BLUE).LogMsg(logger.LOG_INFO, "RENDER", string(str))
		if strings.Contains(string(str), "100%"){
			completed = true
			mainLog.LogMsg(logger.LOG_INFO, "RENDER", "RENDER COMPLETED")
			break
		}
	}
	if completed{
		pd.state = STATE_READY
		packet := &PacketFrameCompleted{LargePacketData:LargePacketData{PacketData: PacketData{PACKET_FRAMECOMPLETED, pd.client}}}
		packet.ProjectName = pd.projectName
		packet.FrameId = frameid
		var path string
		files, _ := ioutil.ReadDir(filepath.Join(rmrfarm.conf.Workspace, "/renderman/", pd.projectName, "images"))
while1:
		for _, file := range files{
			spl := strings.Split(file.Name(), ".")
			for _, substr := range spl {
				if i, err := strconv.Atoi(substr); err == nil && int32(i) == frameid{
					path = filepath.Join(rmrfarm.conf.Workspace, "/renderman/", pd.projectName, "images", file.Name())
					break while1
				}
			}
		}

		mainLog.LogMsg(logger.LOG_INFO, "RENDER","frame rendered  :", path)
		packet.Filepath = path
		rmrfarm.masterHandler.linker.SendPacket(packet)
		pd.changeState(STATE_READY)
	}
}

func (pd *projectData) checkAndRequestRequiredFile() {
	mainLog.SetColor(logger.COLOR_GREEN).LogMsg(logger.LOG_INFO, "PROJECT", "CHECKING REQUIRED FILES")
	fileData := pd.checkRequiredFile()
	if len(fileData) > 0 {
		for _, f := range fileData{
			mainLog.SetColor(logger.COLOR_GREEN).LogMsg(logger.LOG_INFO, "PROJECT", filepath.Join(f.Path, f.File), "NOT EXISTING, ADDING TO REQUEST")
		}
		rmrfarm.masterHandler.linker.SendPacket(&PacketRequestFile{PacketData{PACKET_REQUESTFILE, pd.client}, fileData})
	} else {
		mainLog.SetColor(logger.COLOR_GREEN).LogMsg(logger.LOG_INFO, "PROJECT", "EVERY FILE HAVE BEEN RECEIVED")
		pd.changeState(STATE_READY)
	}
}

func (pd *projectData) checkRequiredFile() []FileData {
	var fileToRequest []FileData
	for id, file := range pd.fileList {
		if !file.IsRequested {
			if _, err := os.Stat(filepath.Join(rmrfarm.conf.Workspace, strings.Replace(file.Path, ":", "_", -1), file.File)); os.IsNotExist(err) {
				mainLog.SetColor(logger.COLOR_GREEN).LogMsg(logger.LOG_INFO, "PROJECT", "REQUIRED FILE ADDED TO REQUEST :", file.File)
				fileToRequest = append(fileToRequest, file)
				pd.fileList[id].IsRequested = true
			}
		}
	}
	return fileToRequest
}

func (pd *projectData) CompatibiliseRIB(){
	mainLog.SetColor(logger.COLOR_BLUE).LogMsg(logger.LOG_INFO, "PROJECT", "CONVERTING RIB PATH FOR COMPATIBILITY WITH THIS PC")
	filepath.Walk(filepath.Join(rmrfarm.conf.Workspace, "renderman", pd.projectName), pd.replaceAbsPath)
	mainLog.SetColor(logger.COLOR_BLUE).LogMsg(logger.LOG_INFO, "PROJECT", "COMPLETED")
}


func  (pd *projectData) replaceAbsPath(path string, info os.FileInfo, err error) error{
	filename := filepath.Base(info.Name())
	if !strings.Contains(filename, ".rib"){
		return nil
	}
	mainLog.SetColor(logger.COLOR_GREEN).LogMsg(logger.LOG_INFO, "PROJECT", "Will convert rib for local path", path)

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
