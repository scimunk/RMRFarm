package main

const (
	FRAMESTATE_WAITING = iota
	FRAMESTATE_RENDERING
	FRAMESTATE_PAUSED
	FRAMESTATE_COMPLETED
)

type frame struct {
	frameId int32 `yaml:"frameId"`
	renderedBy string `yaml:"renderedBy"`
	state   int8 `yaml:"state"`
}
