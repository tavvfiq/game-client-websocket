package game

import (
	"encoding/json"
	"fmt"
	"game-client-websocket/internal/pkg/model"
	"game-client-websocket/internal/pkg/networking"
	"log"
	"net/url"

	"github.com/JoelOtter/termloop"
	"github.com/gorilla/websocket"
	"github.com/tavvfiq/game-server-websocket/pkg/event"
	smdl "github.com/tavvfiq/game-server-websocket/pkg/model"
)

type player struct {
	*termloop.Entity
	ID         string
	networking *networking.Network
}

func newPlayer(id string, networking *networking.Network) *player {
	entity := termloop.NewEntity(0, 0, 1, 1)
	return &player{Entity: entity, ID: id, networking: networking}
}

func (p *player) Tick(tevent termloop.Event) {
	if tevent.Type == termloop.EventKey { // Is it a keyboard event?
		oplayer := model.OtherPlayer{
			ID: p.ID,
		}
		x, y := p.Position()
		switch tevent.Key { // If so, switch on the pressed key.
		case termloop.KeyArrowRight:
			oplayer.State.DeltaX = 1
			p.SetPosition(x+1, y)
		case termloop.KeyArrowLeft:
			oplayer.State.DeltaX = -1
			p.SetPosition(x-1, y)
		case termloop.KeyArrowUp:
			oplayer.State.DeltaY = -1
			p.SetPosition(x, y-1)
		case termloop.KeyArrowDown:
			oplayer.State.DeltaY = 1
			p.SetPosition(x, y+1)
		}
		b, _ := json.Marshal(oplayer)
		payload := smdl.SocketPayload{
			EventType: event.STATE_UPDATE,
			Data:      b,
		}
		err := p.networking.Conn.WriteJSON(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

type Game struct {
	serverID   string
	player     *player
	engine     Engine
	networking *networking.Network
	state      State
}

var addr = "localhost:8080"

func New(serverID, playerID string) *Game {
	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws", RawQuery: fmt.Sprintf("serverID=%s&playerID=%s", serverID, playerID)}
	networking, err := networking.NewConnection(&u)
	if err != nil {
		panic(err)
	}
	state := NewState()
	g := &Game{
		engine:     NewEngine(),
		state:      state,
		networking: networking,
	}
	g.serverID = serverID
	g.player = newPlayer(playerID, networking)
	g.engine.Init(nil)
	return g
}

func (g *Game) Start() {
	g.engine.AddNewEntity(g.player)
	go func() {
		for {
			resp := smdl.SocketResponse{}
			err := g.networking.ReadJSON(&resp)
			if err != nil {
				log.Println("read:", err)
				return
			}
			switch resp.EventType {
			case event.NEW_CONNECTION:
				newPlayer := model.NewOtherPlayer(nil)
				err := json.Unmarshal(resp.Data, &newPlayer)
				if err != nil {
					// less likely will be error
					log.Println(err)
				}
				g.state.AddPlayer(newPlayer)
				g.engine.AddNewEntity(newPlayer)
			case event.PLAYER_DISCONNECT:
				p := model.OtherPlayer{}
				err := json.Unmarshal(resp.Data, &p)
				if err != nil {
					// less likely will be error
					log.Println(err)
				}
				player := g.state.GetPlayer(p.ID)
				g.engine.RemoveEntity(player)
				g.state.RemovePlayer(p.ID)
			case event.SYNC_STATE:
				pls := Players{}
				err := json.Unmarshal(resp.Data, &pls)
				if err != nil {
					// less likely will be error
					panic(err)
				}
				for _, p := range pls {
					pState := &model.PlayerState{X: p.State.X, Y: p.State.Y}
					newP := model.NewOtherPlayer(pState)
					newP.ID = p.ID
					newP.State = p.State
					pls[p.ID] = newP
					if p.ID != g.player.ID {
						g.engine.AddNewEntity(newP)
					}
				}
				g.state.AddPlayers(pls)
			case event.STATE_UPDATE:
				player := model.NewOtherPlayer(nil)
				err := json.Unmarshal(resp.Data, &player)
				if err != nil {
					// less likely will be error
					log.Println(err)
				}
				defer g.state.UpdatePlayer(player)
				p := g.state.GetPlayer(player.ID)
				p.State = player.State
				g.engine.UpdateEntity(p, int(p.State.X), int(p.State.Y))
			}
		}
	}()
	// register the player
	oplayer := model.OtherPlayer{
		ID: g.player.ID,
	}
	b, _ := json.Marshal(oplayer)
	payload := smdl.SocketPayload{
		EventType: event.NEW_CONNECTION,
		Data:      b,
	}
	g.networking.Conn.WriteJSON(payload)
	g.engine.Start()
	err := g.networking.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
		return
	}
}
