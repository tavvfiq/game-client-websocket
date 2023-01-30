package game

import (
	"game-client-websocket/internal/pkg/model"
	"sync"
)

type Players map[string]*model.OtherPlayer

type State interface {
	GetPlayers() Players
	GetPlayer(playerID string) *model.OtherPlayer
	RemovePlayer(playerID string)
	UpdatePlayer(player *model.OtherPlayer)
	AddPlayer(p *model.OtherPlayer)
	AddPlayers(p Players)
}

type _state struct {
	mtx     sync.Mutex
	Players Players
}

func NewState() State {
	return &_state{
		Players: Players{},
	}
}

func (s *_state) AddPlayer(p *model.OtherPlayer) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.Players[p.ID] = p
}

func (s *_state) AddPlayers(ps Players) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.Players = ps
}

func (s *_state) GetPlayers() Players {
	return s.Players
}

func (s *_state) GetPlayer(playerID string) *model.OtherPlayer {
	return s.Players[playerID]
}

func (s *_state) RemovePlayer(playerID string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	newPlayers := Players{}
	for _, p := range s.Players {
		if p.ID != playerID {
			newPlayers[p.ID] = p
		}
	}
	s.Players = newPlayers
}

func (s *_state) UpdatePlayer(player *model.OtherPlayer) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.Players[player.ID] = player
}
