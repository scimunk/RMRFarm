package main

import (
	"bufio"
	"fmt"
	"gopkg.in/v1/yaml"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type config struct {
	Workspace string `yaml:"Workspace"`
	SlaveName string `yaml:"SlaveName"`
	Ip        string `yaml:"IP"`
}

func loadConfiguration() *config {
	file, err := ioutil.ReadFile("slaveConfig.cfg")
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

	config.Workspace = getInput("Define slave workspace path :", exists)
	config.SlaveName = getInput("Define slave Name :", isString)
	config.Ip = getInput("define Slave IP :", checkIp)

	fmt.Println("Writing configuration")
	f, err := os.Create("slaveConfig.cfg")
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
	err := os.MkdirAll(filepath.Clean(path), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func isString(str string) error{
	return nil
}