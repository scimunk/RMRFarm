package main

import "github.com/epixerion/RMRFarm/logger"

type frameManager struct {
	FrameList map[int32]*frame `yaml:"renderedFrame"`
}

func newFrameManager() *frameManager {
	frameManager := &frameManager{}
	frameManager.FrameList = make(map[int32]*frame)
	return frameManager
}

func (fm *frameManager) addFrame(frameId int32){
	mainLog.SetColor(logger.COLOR_BLUE).LogMsg(logger.LOG_INFO, "PROJECT", "frame added :", frameId)
	fm.FrameList[frameId] = &frame{frameId, "", FRAMESTATE_WAITING}
}

func (fM *frameManager) updateFrameManager() {

}

func (fm *frameManager) GetFrameToRender() *frame{
	for _, frame := range fm.FrameList{
		if frame.state == FRAMESTATE_WAITING {
			return frame
		}
	}
	return nil
}