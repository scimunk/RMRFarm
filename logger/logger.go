package logger

import (
	"fmt"
	"os"
	"github.com/shiena/ansicolor"
	"io"
	"strings"
)

const (
	LOG_INFO = iota
	LOG_HIGH
	LOG_MEDIUM
	LOG_LOW
)

type Logger interface{
	LogMsg(int8,string,...interface{})
	LogErr(int8,string,...interface{})
	LogWarn(int8,string,...interface{})
	SetColor(string) Logger
}

type LoggerData struct {
	logLevel int8
	msg io.Writer
	warn io.Writer
	err io.Writer
	nextColor string
}

func NewLogger(logLevel int8) Logger{
	logger := &LoggerData{}
	logger.logLevel = logLevel
	logger.msg = ansicolor.NewAnsiColorWriter(os.Stdout)
	logger.warn = ansicolor.NewAnsiColorWriter(os.Stdout)
	logger.err = ansicolor.NewAnsiColorWriter(os.Stdout)
	return logger
}

func (l *LoggerData) SetColor(colorCode string) Logger{
	l.nextColor = colorCode
	return l
}

func (l *LoggerData) LogMsg(level int8,cat string, data ...interface{}){
	if level <= l.logLevel{
		fmt.Fprintln(l.msg,l.nextColor, "MSG ["+cat+ "] :",strings.Trim(fmt.Sprint(data), "[]"))
		l.nextColor = "\x1b[0m"
	}
}

func (l *LoggerData) LogErr(level int8, cat string,data ...interface{}){
	if level <= l.logLevel{
		fmt.Fprintln(l.warn,l.nextColor, "ERROR[",cat, "] :",strings.Trim(fmt.Sprint(data), "[]"))
		l.nextColor = "\x1b[0m"
	}
}

func (l *LoggerData) LogWarn(level int8, cat string,data ...interface{}){
	if level <= l.logLevel{
		fmt.Fprintln(l.warn, l.nextColor, "WARN[",cat, "] :",strings.Trim(fmt.Sprint(data), "[]"))
		l.nextColor = "\x1b[0m"
	}
}