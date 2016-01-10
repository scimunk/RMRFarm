package main

type frameManager struct {
	RenderedFrame []frame `yaml:"renderedFrame"`
	FrameToRender []frame `yaml:"frameToRender"`
}

func newFrameManager() *frameManager {
	frameManager := &frameManager{}
	return frameManager
}

func (fM *frameManager) updateFrameManager() {

}
