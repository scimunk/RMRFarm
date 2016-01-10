package main

import (
	"bufio"
	"github.com/epixerion/RMRFarm/logger"
	"os"
	"strconv"
	"strings"
)

var commandList map[string]func([]string)

func writeCommand() {
	commandList = make(map[string]func([]string))
	commandList["stats"] = commandStats
	commandList["help"] = commandHelp
	commandList["project"] = commandProject
	commandList["exit"] = commandExit
}

func handleCommand() {
	writeCommand()
	reader := bufio.NewReader(os.Stdin)
	mainLog.LogMsg(logger.LOG_INFO, "INFO", "Hello, Welcome to RMRFarm Manager ! type help for display list of usefull command")
	for {
		input, _ := reader.ReadString('\n')
		data := strings.Split(strings.Trim(input, "\r\n"), " ")
		if commandList[data[0]] != nil {
			commandList[data[0]](data[1:len(data)])
		}
	}
}

func commandCheck(cmd []string, cmdrequired int) bool {
	if len(cmd) < cmdrequired {
		mainLog.LogWarn(logger.LOG_INFO, "CMD", "not enough argument")
		return false
	}
	return true
}

func commandProject(cmd []string) {
	if len(cmd) <= 0 {
		return
	}
	switch cmd[0] {
	case "new":
		if !commandCheck(cmd, 2) {
			return
		}
		if id, err := strconv.Atoi(cmd[1]); id >= 0 && err == nil {
			rmrfarm.projectManager.newProjectId(id)
		}
	case "start":
		rmrfarm.projectManager.projectList[0].startRenderProject()
	case "list":
		mainLog.SetColor("\x1b[33m").LogMsg(logger.LOG_INFO, "PROJECT", "Availlable Project :")
		i := 0
		for _, project := range rmrfarm.projectManager.getAvaillableProject() {
			mainLog.SetColor("\x1b[34m").LogMsg(logger.LOG_INFO,"PROJECT", strconv.Itoa(i) + ") " + project)
			i++
		}
	case "help":
		mainLog.SetColor("\x1b[33m").LogMsg(logger.LOG_INFO,"HELP","Availlable Command for project")
		mainLog.SetColor("\x1b[34m").LogMsg(logger.LOG_INFO,"HELP","new [id]	Create new project (see project list for availlable project)")
		mainLog.SetColor("\x1b[34m").LogMsg(logger.LOG_INFO,"HELP","start [id]	Start Project Rendering")
		mainLog.SetColor("\x1b[34m").LogMsg(logger.LOG_INFO,"HELP","list	List availlable project for creation")
	}
}

func commandHelp(cmd []string) {
	mainLog.SetColor("\x1b[33m").LogMsg(logger.LOG_INFO,"HELP","Availlable Command")
	mainLog.SetColor("\x1b[34m").LogMsg(logger.LOG_INFO,"HELP","stats")
	mainLog.SetColor("\x1b[34m").LogMsg(logger.LOG_INFO,"HELP","project")
	mainLog.SetColor("\x1b[33m").LogMsg(logger.LOG_INFO,"HELP","Type \"[command] help\" for information about the command")
}

func commandStats(cmd []string) {
	if len(cmd) <= 0 {
		return
	}
	switch cmd[0] {
	case "help":
		mainLog.SetColor("\x1b[33m").LogMsg(logger.LOG_INFO,"HELP",len(rmrfarm.slaveManager.slaveData) , " RenderSlave are availlable for rendering")
	}
}

func commandExit(cmd []string){
	exit = false
}