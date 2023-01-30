package networking

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type Network struct {
	*websocket.Conn
}

func NewConnection(u *url.URL) (*Network, error) {
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	return &Network{
		Conn: c,
	}, nil
}

func (n *Network) Close() error {
	return n.Conn.Close()
}
