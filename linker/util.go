package linker

import (
	"github.com/epixerion/RMRFarm/logger"
)

func logMsg(logger logger.Logger, level int8,cat string,msg string,data ...interface{}){
	if logger == nil{
		return
	}
	logger.LogMsg(level, cat,msg, data)
}

func logErr(logger logger.Logger,level int8, cat string,msg string,data ...interface{}){
	if logger == nil{
		return
	}
	logger.LogErr(level, cat,msg, data)
}

func logWarn(logger logger.Logger,level int8, cat string,msg string, data ...interface{}){
	if logger == nil{
		return
	}
	logger.LogWarn(level, cat, msg,data)
}