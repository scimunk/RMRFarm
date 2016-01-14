package main

import (
	"github.com/epixerion/RMRFarm/linker"
	. "github.com/epixerion/RMRFarm/rmrfarm"
	"os"
	"path/filepath"
	"strings"
	"github.com/epixerion/RMRFarm/logger"
)

type projectManager struct {
	projectHook *projectData
	projectList    []*projectData
}

func startProjectManager() *projectManager {
	projectManager := &projectManager{}
	return projectManager
}

func (pm *projectManager) startProject(data *PacketNewProject, client linker.Client) {
	project := newProject(data, client)
	pm.projectList = append(pm.projectList, project)
}

func (pm *projectManager) updateProjectManager() {
	for _, project := range pm.projectList {
		project.updateProject()
	}
}

func (pm *projectManager) HandleSendFile(packet *PacketSendFile) {
	err := os.MkdirAll(filepath.Join(rmrfarm.conf.Workspace, strings.Replace(packet.Path, ":", "_", -1)), os.ModePerm)
	if err != nil {
		mainLog.SetColor(logger.COLOR_RED).LogErr(logger.LOG_INFO, "PROJECT", "Could'nt not create directory :", err)
	}
	if packet.IsExterne {
		err = os.Rename(packet.Filepath, filepath.Join(rmrfarm.conf.Workspace, strings.Replace(packet.Path, ":", "_", -1), packet.FileName))
	}else {
		err = os.Rename(packet.Filepath, filepath.Join(rmrfarm.conf.Workspace, strings.Replace(packet.Path, ":", "_", -1), packet.FileName))
	}
	if err != nil {
		mainLog.SetColor(logger.COLOR_RED).LogErr(logger.LOG_INFO, "PROJECT", "Couldn't Move Temp File", err)
	}
	mainLog.SetColor(logger.COLOR_GREEN).LogMsg(logger.LOG_INFO, "PROJECT", "FILE ",filepath.Join(rmrfarm.conf.Workspace, strings.Replace(packet.Path, ":", "_", -1), packet.FileName), "Have Been Received")
}

func (pm *projectManager) renderFrame(packet *PacketRenderFrame){
	for _ , project := range pm.projectList{
		if project.projectName == packet.ProjectName{
			project.startRender(packet.Camera, packet.FrameId)
			break
		}
	}
}