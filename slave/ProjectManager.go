package main

import (
	"github.com/epixerion/RMRFarm/linker"
	. "github.com/epixerion/RMRFarm/rmrfarm"
	"fmt"
)

type projectManager struct {
	currentProject *projectData
	projectList    []*projectData
}

func startProjectManager() *projectManager {
	projectManager := &projectManager{}
	return projectManager
}

func (pm *projectManager) startProject(data *PacketNewProject, client linker.Client) {
	pm.projectList = append(pm.projectList, newProject(data, client))
	pm.currentProject = pm.projectList[0]
	pm.currentProject.generateProject()
}

func (pm *projectManager) updateProjectManager() {
	for _, project := range pm.projectList {
		project.updateProject()
	}
}

func (pm *projectManager) UpdateSlaveReadyness(){
	ready := []string{}
	fmt.Println("updating slave readyness")
	for _, project := range pm.projectList{
		if project.state == STATE_READY{
			ready = append(ready, project.projectName)
		}
	}
	pm.currentProject.client.GetConn().SendPacket(&PacketSlaveReady{PacketData{PACKET_SLAVEREADY, nil}, ready, pm.currentProject.state == STATE_READY})
}

func (pm *projectManager) renderFrame(packet *PacketRenderFrame){
	for _ , project := range pm.projectList{
		if project.projectName == packet.ProjectName{
			project.startRender(packet.Camera, packet.FrameId)
			break
		}
	}
}