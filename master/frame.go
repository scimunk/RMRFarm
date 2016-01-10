package main

const (
	FRAMESTATE_WAITING = iota
	FRAMESTATE_RENDERING
	FRAMESTATE_PAUSED
	FRAMESTATE_COMPLETED
)

type frame struct {
	frameId int32
	state   int8
}
