package model

import "github.com/JoelOtter/termloop"

type Entity interface {
	Draw(s *termloop.Screen)
	Tick(ev termloop.Event)
	Position() (int, int)
	Size() (int, int)
	SetCell(x, y int, c *termloop.Cell)
	Fill(c *termloop.Cell)
	ApplyCanvas(c *termloop.Canvas)
	SetCanvas(c *termloop.Canvas)
	SetPosition(x, y int)
}
