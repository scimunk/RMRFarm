package main

import (
	"github.com/epixerion/RMRFarm/linker"
	. "github.com/epixerion/RMRFarm/rmrfarm"
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
}

func (pm *projectManager) updateProjectManager() {
	for _, project := range pm.projectList {
		project.updateProject()
	}
}
