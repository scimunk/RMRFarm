package main

import (
	"bufio"
	"github.com/epixerion/RMRFarm/logger"
	"os"
	"strconv"
	"strings"
)

var commandList map[string]func([]string)
var reader *bufio.Reader

func writeCommand() {
	commandList = make(map[string]func([]string))
	commandList["stats"] = commandStats
	commandList["help"] = commandHelp
	commandList["project"] = commandProject
	commandList["renderslave"] = commandRenderSlave
	commandList["exit"] = commandExit
}

func handleCommand() {
	writeCommand()
	reader = bufio.NewReader(os.Stdin)
	mainLog.SetColor(logger.COLOR_GREEN).LogMsg(logger.LOG_INFO, "INFO", "Hello, Welcome to RMRFarm Manager ! type help for display list of usefull command")
	for {
		data := readCmd()
		if commandList[data[0]] != nil {
			commandList[data[0]](data[1:len(data)])
		}
	}
}

func readCmd() []string{
	input, _ := reader.ReadString('\n')
	data := strings.Split(strings.Trim(input, "\r\n"), " ")
	return data
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
		mainLog.SetColor(logger.COLOR_YELLOW).LogMsg(logger.LOG_INFO, "PROJECT", "Availlable Project :")
		i := 1
		mainLog.SetColor(logger.COLOR_RED).LogMsg(logger.LOG_INFO,"PROJECT", "0) Cancel")
		for _, project := range rmrfarm.projectManager.getAvaillableProject() {
			mainLog.SetColor(logger.COLOR_GREEN).LogMsg(logger.LOG_INFO,"PROJECT", strconv.Itoa(i) + ") " + project)
			i++
		}
		mainLog.SetColor(logger.COLOR_CYAN).LogMsg(logger.LOG_INFO, "Type the id of the project you want to start :")
		cmd := readCmd()
		if id, err := strconv.Atoi(cmd[0]); id > 0 && err == nil {
			rmrfarm.projectManager.newProjectId(id-1)
		}
	case "start":
		mainLog.SetColor(logger.COLOR_CYAN).LogMsg(logger.LOG_INFO, "PROJECT", "Choose Project to start :")
		i := 0
		for _, project := range rmrfarm.projectManager.getGeneratedProject() {
			mainLog.SetColor(logger.COLOR_GREEN).LogMsg(logger.LOG_INFO,"PROJECT", strconv.Itoa(i)+ ") "+project)
			i++
		}
		cmd := readCmd()
		if id, err := strconv.Atoi(cmd[0]); id >= 0 && id < len(rmrfarm.projectManager.ProjectList) && err == nil {
			rmrfarm.projectManager.ProjectList[id].startRenderProject()
		}else{
			mainLog.SetColor(logger.COLOR_RED).LogMsg(logger.LOG_INFO, "PROJECT", "Cancelled")
		}
	case "list":
		mainLog.SetColor(logger.COLOR_YELLOW).LogMsg(logger.LOG_INFO, "PROJECT", "Generated Project :")
		i := 0
		for _, project := range rmrfarm.projectManager.getGeneratedProject() {
			mainLog.SetColor(logger.COLOR_GREEN).LogMsg(logger.LOG_INFO,"PROJECT", project)
			i++
		}
	case "help":
		mainLog.SetColor(logger.COLOR_YELLOW).LogMsg(logger.LOG_INFO,"HELP","Availlable Command for project")
		mainLog.SetColor(logger.COLOR_BLUE).LogMsg(logger.LOG_INFO,"HELP","new    Create new project (see project list for availlable project)")
		mainLog.SetColor(logger.COLOR_BLUE).LogMsg(logger.LOG_INFO,"HELP","start   Start Project Rendering")
		mainLog.SetColor(logger.COLOR_BLUE).LogMsg(logger.LOG_INFO,"HELP","list	List availlable project for creation")
	}
}

func commandHelp(cmd []string) {
	mainLog.SetColor(logger.COLOR_YELLOW).LogMsg(logger.LOG_INFO,"HELP","Availlable Command")
	mainLog.SetColor(logger.COLOR_BLUE).LogMsg(logger.LOG_INFO,"HELP","stats")
	mainLog.SetColor(logger.COLOR_BLUE).LogMsg(logger.LOG_INFO,"HELP","project")
	mainLog.SetColor(logger.COLOR_BLUE).LogMsg(logger.LOG_INFO,"HELP","renderslave")
	mainLog.SetColor(logger.COLOR_BLUE).LogMsg(logger.LOG_INFO,"HELP","exit")
	mainLog.SetColor(logger.COLOR_YELLOW).LogMsg(logger.LOG_INFO,"HELP","Type \"[command] help\" for information about the command")
}

func commandStats(cmd []string) {
	if len(cmd) <= 0 {
		return
	}
	switch cmd[0] {
		default:
		mainLog.SetColor(logger.COLOR_YELLOW).LogMsg(logger.LOG_INFO,"HELP",len(rmrfarm.slaveManager.slaveData) , " RenderSlave are availlable for rendering")
	}
}

func commandRenderSlave(cmd []string) {
	if len(cmd) <= 0 {
		return
	}
	switch cmd[0] {
	case "status":
		if len(rmrfarm.slaveManager.slaveData) == 0{
			mainLog.SetColor(logger.COLOR_YELLOW).LogMsg(logger.LOG_INFO,"HELP","No renderslave availlable")
		}
		for _, slave := range rmrfarm.slaveManager.slaveData {
			slave.printSlaveInfo()
		}
	case "help":
		mainLog.SetColor(logger.COLOR_YELLOW).LogMsg(logger.LOG_INFO,"HELP",len(rmrfarm.slaveManager.slaveData) , " RenderSlave are availlable for rendering")
	default:
		mainLog.SetColor(logger.COLOR_YELLOW).LogMsg(logger.LOG_INFO,"HELP", "renderslave help")
	}
}

func commandExit(cmd []string){
	exit = false
}