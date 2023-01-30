package game

import (
	"game-client-websocket/internal/pkg/model"

	"github.com/JoelOtter/termloop"
)

type Engine interface {
	Start()
	AddNewEntity(entity model.Entity)
	RemoveEntity(entity model.Entity)
	Init(level termloop.Level)
	UpdateEntity(entity model.Entity, x, y int)
}

type _engine struct {
	tloop *termloop.Game
	level termloop.Level
}

func NewEngine() Engine {
	return &_engine{
		tloop: termloop.NewGame(),
	}
}

func (e *_engine) Init(level termloop.Level) {
	if level == nil {
		level = termloop.NewBaseLevel(termloop.Cell{
			Bg: termloop.ColorGreen,
			Fg: termloop.ColorBlack,
			Ch: 'v',
		})
	}
	e.level = level
	e.level.AddEntity(termloop.NewRectangle(10, 10, 50, 20, termloop.ColorBlue))
}

func (e *_engine) Start() {
	e.setLevel()
	e.tloop.Start()
}

func (e *_engine) setLevel() {
	if e.level != nil {
		e.tloop.Screen().SetLevel(e.level)
	}
}

func (e *_engine) AddNewEntity(entity model.Entity) {
	entity.SetCell(0, 0, &termloop.Cell{Fg: termloop.ColorRed, Ch: '옷'})
	e.level.AddEntity(entity)
}

func (e *_engine) RemoveEntity(entity model.Entity) {
	e.level.RemoveEntity(entity)
}

func (e *_engine) UpdateEntity(entity model.Entity, x, y int) {
	entity.SetPosition(x, y)
}
