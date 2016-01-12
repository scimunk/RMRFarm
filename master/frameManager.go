package main

import "github.com/epixerion/RMRFarm/logger"

type frameManager struct {
	RenderedFrame map[int32]*frame `yaml:"renderedFrame"`
	RenderingFrame map[int32]*frame `yaml:"RenderingFrame"`
	FrameToRender map[int32]*frame `yaml:"frameToRender"`
}

func newFrameManager() *frameManager {
	frameManager := &frameManager{}
	frameManager.FrameToRender = make(map[int32]*frame)
	frameManager.RenderingFrame = make(map[int32]*frame)
	frameManager.RenderedFrame = make(map[int32]*frame)
	return frameManager
}

func (fm *frameManager) addFrame(frameId int32){
	mainLog.SetColor(logger.COLOR_BLUE).LogMsg(logger.LOG_INFO, "PROJECT", "frame added :", frameId)
	fm.FrameToRender[frameId] = &frame{frameId, FRAMESTATE_WAITING}
}

func (fM *frameManager) updateFrameManager() {

}

func (fm *frameManager) RenderFrame() *frame{
	for id, frame := range fm.FrameToRender{
		delete(fm.FrameToRender, id)
		fm.RenderingFrame[id] = frame
		return frame
	}
	return nil
}