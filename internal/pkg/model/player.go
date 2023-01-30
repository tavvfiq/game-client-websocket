package model

import (
	"github.com/JoelOtter/termloop"
)

type OtherPlayer struct {
	*termloop.Entity
	ID    string
	State PlayerState
}

func NewOtherPlayer(state *PlayerState) *OtherPlayer {
	entity := termloop.NewEntity(0, 0, 1, 1)
	if state != nil {
		entity.SetPosition(int(state.X), int(state.Y))
	}
	return &OtherPlayer{
		Entity: entity,
	}
}

type PlayerState struct {
	X      float64
	Y      float64
	DeltaX float64
	DeltaY float64
}
