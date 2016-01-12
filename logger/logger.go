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

const (
	TEXT_BOLD = iota
	TEXT_UNDERLINE
	TEXT_BLINK

	COLOR_BLACK
	COLOR_RED
	COLOR_GREEN
	COLOR_YELLOW
	COLOR_BLUE
	COLOR_MAGENTA
	COLOR_CYAN
	COLOR_WHITE
	COLOR_LIGHTGRAY
	COLOR_LIGHTRED
	COLOR_LIGHTGREEN
	COLOR_LIGHTYELLOW
	COLOR_LIGHTBLUE
	COLOR_LIGHTMAGENTA
	COLOR_LIGHTCYAN
	COLOR_LIGHTWHITE

	BACKGROUND_BLACK
	BACKGROUND_RED
	BACKGROUND_GREEN
	BACKGROUND_YELLOW
	BACKGROUND_BLUE
	BACKGROUND_MAGENTA
	BACKGROUND_CYAN
	BACKGROUND_WHITE
	BACKGROUND_LIGHTGRAY
	BACKGROUND_LIGHTRED
	BACKGROUND_LIGHTGREEN
	BACKGROUND_LIGHTYELLOW
	BACKGROUND_LIGHTBLUE
	BACKGROUND_LIGHTMAGENTA
	BACKGROUND_LIGHTCYAN
	BACKGROUND_LIGHTWHITE
)

var colorCodeMap = []string{
	"\x1b[1m",
	"\x1b[4m",
	"\x1b[5m",

	"\x1b[30m",
	"\x1b[31m",
	"\x1b[32m",
	"\x1b[33m",
	"\x1b[34m",
	"\x1b[35m",
	"\x1b[36m",
	"\x1b[37m",

	"\x1b[90m",
	"\x1b[91m",
	"\x1b[92m",
	"\x1b[93m",
	"\x1b[94m",
	"\x1b[95m",
	"\x1b[96m",
	"\x1b[97m",

	"\x1b[40m",
	"\x1b[41m",
	"\x1b[42m",
	"\x1b[43m",
	"\x1b[44m",
	"\x1b[45m",
	"\x1b[46m",
	"\x1b[47m",

	"\x1b[100m",
	"\x1b[101m",
	"\x1b[102m",
	"\x1b[103m",
	"\x1b[104m",
	"\x1b[105m",
	"\x1b[106m",
	"\x1b[107m",
}

type Logger interface{
	LogMsg(int8,string,...interface{})
	LogErr(int8,string,...interface{})
	LogWarn(int8,string,...interface{})
	SetColor(...int) Logger
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

func (l *LoggerData) SetColor(colorCode ...int) Logger{
	if len(colorCode) > 0{
		l.nextColor = ""
	}
	for _, color := range colorCode{
		l.nextColor += colorCodeMap[color]
	}
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