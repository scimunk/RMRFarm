package main

import (
	"github.com/epixerion/RMRFarm/logger"
	"gopkg.in/v1/yaml"
	"io/ioutil"
	"path/filepath"
)

type projectManager struct {
	ProjectList []*projectData `yaml:"projectList"`
}

func startProjectManager() *projectManager {
	projectManager := &projectManager{}
	file, err := ioutil.ReadFile("project.cfg")
	if err != nil {
	}
	yaml.Unmarshal(file, &projectManager)
	return projectManager
}

func (pm *projectManager) updateProjectManager(){
	for _, project := range pm.ProjectList{
		project.updateProject()
	}
}

func (pm *projectManager) newProjectId(id int) {
	if len(pm.getAvaillableProject()) > id {
		projectName := pm.getAvaillableProject()[id]
		project := newProject(projectName)
		pm.ProjectList = append(pm.ProjectList, project)
		pm.saveProject()
	} else {
		mainLog.SetColor(logger.COLOR_RED).LogWarn(1, "PROJECT", "project id", id, "do not exist")
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

func (pm *projectManager) getGeneratedProject() []string{
	str := []string{}
	for _, p := range pm.ProjectList{
		str = append(str, p.ProjectName)
	}
	return str
}

func (pm *projectManager) saveProject() {
	mainLog.SetColor(logger.COLOR_GREEN).LogMsg(logger.LOG_INFO, "PROJECTMANAGER", "SAVING PROJECT DATA")
	data, err := yaml.Marshal(pm)
	ioutil.WriteFile("project.cfg", data, 0644)
	check(err)
	mainLog.SetColor(logger.COLOR_GREEN).LogMsg(logger.LOG_INFO, "PROJECTMANAGER", "SAVED!")
}
