package main

import (
	"bufio"
	"errors"
	"fmt"
	"gopkg.in/v1/yaml"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type config struct {
	MayaWorkspace string   `yaml:"mayaWorkspace"`
	SlaveListIp   []string `yaml:"slaveIP"`
}

func loadConfiguration() *config {
	file, err := ioutil.ReadFile("masterConfig.cfg")
	if err != nil {
		return createConf()
	}
	conf := &config{}
	yaml.Unmarshal(file, &conf)
	return conf
}

func createConf() *config {
	fmt.Println("Configuration files not found, automatic creation :")
	config := &config{}

	config.MayaWorkspace = getInput("Define maya workspace path :", exists)
	config.SlaveListIp = append(config.SlaveListIp, getInput("define server :", checkIp))

	fmt.Println("Writing configuration")
	f, err := os.Create("masterConfig.cfg")
	check(err)
	defer f.Close()
	data, err := yaml.Marshal(config)
	check(err)
	fmt.Println("done!")
	f.Write(data)
	return config
}

func getInput(message string, validator func(string) error) string {
	var input string
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(message)
		input, _ = reader.ReadString('\n')
		err := validator((strings.Trim(input, "\r\n")))
		if err != nil {
			fmt.Println(err.Error())
		} else {
			break
		}
	}
	return strings.Trim(input, "\r\n")
}

func checkIp(ip string) error {
	return nil
}

func exists(path string) error {
	_, err := os.Stat(path + "/")
	if err != nil {
		return err
	}
	file, _ := filepath.Glob(path + "/*.mel")
	if file == nil {
		return errors.New("Folder must containt an maya workspace (.mel)")
	}
	return nil
}
