package main

import (
	"archive/zip"
	"fmt"
	. "github.com/epixerion/RMRFarm/rmrfarm"
	"github.com/epixerion/RMRFarm/logger"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"strconv"
)

const (
	PROJECT_STATE_IDLE = iota
	PROJECT_STATE_RENDER
	PROJECT_STATE_FINISHED
)

type projectData struct {
	ProjectName   string        `yaml:"projectName"`
	FileData      []FileData    `yaml:"fileList"`
	FrameManager  *frameManager `yaml:"frameManager"`
	Camera []string `yaml:"CameraList"`
	cameraToRender int
	State int
}

func newProject(projectName string) *projectData {
	projectData := &projectData{ProjectName: projectName}
	projectData.FrameManager = newFrameManager()
	projectData.generateProject()
	projectData.State = PROJECT_STATE_IDLE
	return projectData
}

func (pd *projectData) updateProject(){
	if pd.State == PROJECT_STATE_RENDER{
		if slave := rmrfarm.slaveManager.getSlaveReadyForProject(pd.ProjectName); slave != nil{
			if frame := pd.FrameManager.GetFrameToRender(); frame != nil {
				slave.StartRenderFrame(pd, frame)
				mainLog.SetColor(logger.COLOR_LIGHTGREEN).LogMsg(logger.LOG_INFO, "PROJECT", "Rendering frame : ", frame.frameId)
			}
		}
	}
}

func (pd *projectData) startRenderProject() {
	mainLog.SetColor(logger.COLOR_CYAN).LogMsg(logger.LOG_INFO, "PROJECT", "Select the camera you want to render :")
	for id, camera := range pd.Camera{
		mainLog.SetColor(logger.COLOR_GREEN).LogMsg(logger.LOG_INFO, "PROJECT", id, ")", camera)
	}
	cmd := readCmd()
	if id, err := strconv.Atoi(cmd[0]); id >= 0 && id < len(pd.Camera) && err == nil {
		mainLog.LogMsg(logger.LOG_INFO, "PROJECT", "Starting Rendering of camera", pd.Camera[id])
		pd.cameraToRender = id
	}
	for _, slave := range rmrfarm.slaveManager.getAvaillableSlave(){
		slave.AssignSlaveToProject(pd)
	}
	pd.State = PROJECT_STATE_RENDER
}

func (pd *projectData) generateProject() {
	pd.compressProject()

	mainLog.SetColor(logger.COLOR_MAGENTA).LogMsg(logger.LOG_INFO, "PROJECT", "REGISTERING EXTERNAL ASSET")
	folders, _ := ioutil.ReadDir(filepath.Join(rmrfarm.conf.MayaWorkspace, "/renderman/", pd.ProjectName))
	for _, dir := range folders {
		if dir.IsDir() && strings.HasPrefix(dir.Name(), "rib") {
			frames, _ := ioutil.ReadDir(filepath.Join(rmrfarm.conf.MayaWorkspace, "/renderman/", pd.ProjectName, dir.Name()))
			for _, frame := range frames {
				if id, err := strconv.Atoi(frame.Name()); err == nil && id != 0{
					pd.FrameManager.addFrame(int32(id))
				}
				files, _ := ioutil.ReadDir(filepath.Join(rmrfarm.conf.MayaWorkspace, "/renderman/", pd.ProjectName, dir.Name(), frame.Name()))
				for _, file := range files {
					if !strings.HasSuffix(file.Name(), ".rib") && !strings.HasSuffix(file.Name(), ".rlf"){
						continue
					}
					if strings.Contains(file.Name(), "Shape"){
						pd.addCamera(file.Name()[:strings.Index(file.Name(),"Shape")])
					}
					content, _ := ioutil.ReadFile(filepath.Join(rmrfarm.conf.MayaWorkspace, "/renderman/", pd.ProjectName, dir.Name(), frame.Name(), file.Name()))
					fmt.Println("reading file", filepath.Join(rmrfarm.conf.MayaWorkspace, "/renderman/", pd.ProjectName, dir.Name(), frame.Name(), file.Name()))
					pd.ExtractPath(string(content))
				}
			}
		}
	}
	mainLog.SetColor(logger.COLOR_MAGENTA).LogMsg(logger.LOG_INFO, "PROJECT", "REGISTERING COMPLETED")
	mainLog.SetColor(logger.COLOR_GREEN).LogMsg(logger.LOG_INFO, "PROJECT", "PROJECT", pd.ProjectName, "GENERATED WITH SUCCESS")
}

func (pd *projectData) addCamera(name string){
	for _, cam := range pd.Camera{
		if cam == name{
			return
		}
	}
	mainLog.SetColor(logger.COLOR_CYAN).LogMsg(logger.LOG_INFO, "PROJECT", "Added Camera :", name)
	pd.Camera = append(pd.Camera, name)
}

type microCompress struct {
	writer *zip.Writer
}

func (pd *projectData) compressProject() {
	mainLog.SetColor(logger.COLOR_MAGENTA).LogMsg(logger.LOG_INFO, "PROJECT", "START COMPRESSING PROJECT")
	err := os.MkdirAll(filepath.Join(rmrfarm.conf.MayaWorkspace, "RMRFarm"), os.ModePerm)

	zipPath := filepath.Join(rmrfarm.conf.MayaWorkspace, "RMRFarm", pd.ProjectName +".zip")
	if f, _ := os.Stat(zipPath); f != nil {
		os.Remove(zipPath)
	}

	zipfile, _ := os.Create(zipPath)
	mc := &microCompress{}

	mc.writer = zip.NewWriter(zipfile)
	filepath.Walk(filepath.Join(rmrfarm.conf.MayaWorkspace, "renderman", pd.ProjectName), mc.addToZip)
	err = mc.writer.Close()
	zipfile.Close()
	if err != nil {
		mainLog.LogErr(logger.LOG_HIGH, "COMPRESS", "", err)
	}
	mainLog.SetColor(logger.COLOR_MAGENTA).LogMsg(logger.LOG_INFO, "PROJECT", "COMPRESSION COMPLETED")
}

func (mc *microCompress) addToZip(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}
	newpath, _ := filepath.Rel(rmrfarm.conf.MayaWorkspace, path)
	ftowrite, _ := os.Open(path)
	if strings.Contains(newpath, "images") {
		return nil
	}
	if err != nil {
		fmt.Println(err)
	}
	mainLog.SetColor(logger.COLOR_BLUE).LogMsg(logger.LOG_INFO, "COMPRESS", "Compressing file", newpath)
	f, _ := mc.writer.Create(newpath)
	io.Copy(f, ftowrite)
	ftowrite.Close()
	return nil
}

/*
Extracting all path from the passed file content
 */
func (pd *projectData) ExtractPath(content string) []FileData {
	regex, _ := regexp.Compile("(([A-Z]:\\/)?([a-zA-Z0-9_.\\-]+[\\/\\])+[a-zA-Z0-9_.\\-]*))")
	var pathArray []FileData
	for _, matchedRegex := range regex.FindAllStringSubmatch(string(content), -1) {
		quoted := matchedRegex[0]
		if f, err := os.Stat(quoted); !os.IsNotExist(err) && !f.IsDir() && filepath.IsAbs(quoted){
			fileData := FileData{filepath.Base(quoted), filepath.Dir(quoted), true, false}
			pd.RegisterPath(fileData)
		} else if strings.Index(quoted, "/") != -1 && strings.Index(quoted, "@") == -1 && strings.Index(quoted, ":") == -1 {
			if !strings.HasSuffix(quoted, "/") {
				if strings.Index(quoted, "renderman") == 0 && (strings.Index(quoted, "ribarchive") != -1 || strings.Index(quoted, "textures") != -1) {
					fileData := FileData{filepath.Base(quoted), filepath.Dir(quoted), false, false}
					pd.RegisterPath(fileData)
				}
			}
		}
	}
	return pathArray
}

/*
We register the path in the fileData List, if it not exist
 */
func (pd *projectData) RegisterPath(link FileData) {
	exist := false
	for _, existingcheck := range pd.FileData {
		if existingcheck.File == link.File {
			exist = true
		}
	}
	if !exist {
		var err error
		if link.IsExterne {
			_, err = os.Stat(filepath.Join(link.Path, link.File))
		} else {
			_, err = os.Stat(filepath.Join(rmrfarm.conf.MayaWorkspace, link.Path, link.File))
		}
		if os.IsNotExist(err) {
			mainLog.SetColor(logger.BACKGROUND_RED).LogMsg(logger.LOG_INFO, "DEBUG", err)
			return
		}
		pd.FileData = append(pd.FileData, link)
		mainLog.SetColor(logger.COLOR_BLUE).LogMsg(logger.LOG_INFO, "PROJECT", "Register file", link.Path+link.File, "for remote ressource")
	}
}

func (pd *projectData) handleFrameCompleted(packet *PacketFrameCompleted){
	err := os.MkdirAll(filepath.Join(rmrfarm.conf.MayaWorkspace, "RMRFarm", pd.ProjectName), os.ModePerm)
	if err != nil {
		mainLog.SetColor(logger.COLOR_RED).LogErr(logger.LOG_INFO, "PROJECT", "Could'nt not create directory :", err)
	}
	err = os.Rename(packet.Filepath, filepath.Join(rmrfarm.conf.MayaWorkspace, "RMRFarm", pd.ProjectName, strconv.Itoa(int(packet.FrameId)) + ".exr"))
	if err != nil {
		mainLog.SetColor(logger.COLOR_RED).LogErr(logger.LOG_INFO, "PROJECT", "Couldn't Move Temp File", err)
	}
	mainLog.SetColor(logger.COLOR_GREEN).LogMsg(logger.LOG_INFO, "PROJECT", "FRAME RECEIVED ",
		filepath.Join(rmrfarm.conf.MayaWorkspace, "RMRFarm", pd.ProjectName, strconv.Itoa(int(packet.FrameId)) + ".exr"))
}