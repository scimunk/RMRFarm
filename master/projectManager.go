package main

import (
	"github.com/epixerion/RMRFarm/logger"
	"gopkg.in/v1/yaml"
	"io/ioutil"
	"path/filepath"
)

type projectManager struct {
	projectList []*projectData `yaml:"projectList"`
}

func startProjectManager() *projectManager {
	projectManager := &projectManager{}
	return projectManager
}

func (pm *projectManager) newProjectId(id int) {
	if len(pm.getAvaillableProject()) > id {
		projectName := pm.getAvaillableProject()[id]
		project := newProject(projectName)
		pm.projectList = append(pm.projectList, project)
	} else {
		mainLog.SetColor("\x1b[31m").LogWarn(1, "PROJECT", "project id", id, "do not exist")
	}
}

func (pm *projectManager) getAvaillableProject() []string {
	files, _ := ioutil.ReadDir(filepath.Join(rmrfarm.conf.MayaWorkspace, "renderman"))
	var project []string
	for _, file := range files {
		if file.Name() != "ribarchives" && file.Name() != "textures" {
			project = append(project, file.Name())
		}
	}
	return project
}

func (pm *projectManager) saveProject() {
	mainLog.SetColor("\x1b[32m").LogMsg(logger.LOG_INFO, "PROJECTMANAGER", "SAVING PROJECT DATA")
	data, err := yaml.Marshal(pm)
	ioutil.WriteFile("project.cfg", data, 0644)
	check(err)
	mainLog.SetColor("\x1b[32m").LogMsg(logger.LOG_INFO, "PROJECTMANAGER", "SAVED!")
}
