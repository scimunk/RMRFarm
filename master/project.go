package main

import (
	"archive/zip"
	"fmt"
	. "github.com/epixerion/RMRFarm/rmrfarm"
	"github.com/epixerion/RMRFarm/logger"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type projectData struct {
	projectName   string        `yaml:"projectName"`
	fileData      []FileData    `yaml:"fileList"`
	frameManager  *frameManager `yaml:"frameManager"`
	assignedSlave []string      `yaml:"assignedSlave"`
}

func newProject(projectName string) *projectData {
	projectData := &projectData{projectName: projectName}
	projectData.frameManager = newFrameManager()
	projectData.generateProject()
	return projectData
}

func (pd *projectData) addSlaveToProject(slaveName []string) {

}

func (pd *projectData) startRenderProject() {
	rmrfarm.slaveManager.getAvaillableSlave()[0].AssignSlaveToProject(pd)
	if len(pd.assignedSlave) < 0 {
		mainLog.LogWarn(0, "PROJECT", "No slave assigned to this project, use project slave assign [id] to assign slave eg: project slave assign 5 1 3")
		return
	}

}

func (pd *projectData) generateProject() {
	pd.compressProject()

	mainLog.SetColor("\x1b[35m").LogMsg(logger.LOG_INFO, "PROJECT", "REGISTERING EXTERNAL ASSET")
	folders, _ := ioutil.ReadDir(filepath.Join(rmrfarm.conf.MayaWorkspace, "/renderman/", pd.projectName))
	for _, dir := range folders {
		if dir.IsDir() && strings.HasPrefix(dir.Name(), "rib") {
			frames, _ := ioutil.ReadDir(filepath.Join(rmrfarm.conf.MayaWorkspace, "/renderman/", pd.projectName, dir.Name()))
			for _, frame := range frames {
				files, _ := ioutil.ReadDir(filepath.Join(rmrfarm.conf.MayaWorkspace, "/renderman/", pd.projectName, dir.Name(), frame.Name()))
				for _, file := range files {
					if !strings.HasSuffix(file.Name(), ".rib") {
						continue
					}
					content, _ := ioutil.ReadFile(filepath.Join(rmrfarm.conf.MayaWorkspace, "/renderman/", pd.projectName, dir.Name(), frame.Name(), file.Name()))
					pd.extractRibLink(string(content))
				}
			}
		}
	}
	mainLog.SetColor("\x1b[35m").LogMsg(logger.LOG_INFO, "PROJECT", "REGISTERING COMPLETED")
	mainLog.SetColor("\x1b[32m").LogMsg(logger.LOG_INFO, "PROJECT", "PROJECT", pd.projectName, "GENERATED WITH SUCCESS")
}

type microCompress struct {
	writer *zip.Writer
}

func (pd *projectData) compressProject() {
	mainLog.SetColor("\x1b[35m").LogMsg(logger.LOG_INFO, "PROJECT", "START COMPRESSING PROJECT")
	err := os.MkdirAll(filepath.Join(rmrfarm.conf.MayaWorkspace, "RMRFarm"), os.ModePerm)

	zipPath := filepath.Join(rmrfarm.conf.MayaWorkspace, "RMRFarm", pd.projectName+".zip")
	if f, _ := os.Stat(zipPath); f != nil {
		os.Remove(zipPath)
	}

	zipfile, _ := os.Create(zipPath)
	mc := &microCompress{}

	mc.writer = zip.NewWriter(zipfile)
	filepath.Walk(filepath.Join(rmrfarm.conf.MayaWorkspace, "renderman", pd.projectName), mc.addToZip)
	err = mc.writer.Close()
	zipfile.Close()
	if err != nil {
		mainLog.LogErr(logger.LOG_HIGH, "COMPRESS", "", err)
	}
	mainLog.SetColor("\x1b[35m").LogMsg(logger.LOG_INFO, "PROJECT", "COMPRESSION COMPLETED")
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
	mainLog.SetColor("\x1b[34m").LogMsg(logger.LOG_INFO, "COMPRESS", "Compressing file", newpath)
	f, _ := mc.writer.Create(newpath)
	io.Copy(f, ftowrite)
	ftowrite.Close()
	return nil
}

func (pd *projectData) addUniqueLink(link FileData) {
	exist := false
	for _, existingcheck := range pd.fileData {
		if existingcheck.File == link.File {
			exist = true
		}
	}
	if !exist {
		var err error
		if !link.IsExterne {
			_, err = os.Stat(filepath.Join(link.Path, link.File))
		} else {
			_, err = os.Stat(filepath.Join(rmrfarm.conf.MayaWorkspace, link.Path, link.File))
		}
		if os.IsNotExist(err) {
			return
		}
		pd.fileData = append(pd.fileData, link)
		mainLog.SetColor("\x1b[34m").LogMsg(logger.LOG_INFO, "PROJECT", "Register file", link.Path+link.File, "for remote ressource")
	}
}

func (pd *projectData) extractRibLink(content string) []FileData {
	regex, _ := regexp.Compile("(([A-Z]:\\/)?([a-zA-Z0-9_.\\-]+[\\/\\])+[a-zA-Z0-9_.\\-]*))")
	var linkArray []FileData
	for _, found := range regex.FindAllStringSubmatch(string(content), -1) {
		quoted := found[0]
		if f, err := os.Stat(quoted); !os.IsNotExist(err) && !f.IsDir() {
			fileData := FileData{path.Base(quoted), path.Dir(quoted), true}
			pd.addUniqueLink(fileData)
		} else if strings.Index(quoted, "/") != -1 && strings.Index(quoted, "@") == -1 && strings.Index(quoted, ":") == -1 {
			if !strings.HasSuffix(quoted, "/") {
				if strings.Index(quoted, "ribarchive") != -1 || strings.Index(quoted, "textures") != -1 {
					fileData := FileData{path.Base(quoted), path.Dir(quoted), false}
					pd.addUniqueLink(fileData)
				}
			}
		}
	}
	return linkArray
}
