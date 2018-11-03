package http

import (
	"errors"
	"fmt"
	"log"
	"rockpaperscissor/rockpaperscissor"
	"time"

	"github.com/gorilla/websocket"
)

type hub struct {
	players map[string]*client
	history []rockpaperscissor.Result

	moveC       chan ReadMessage
	resultsC    chan rockpaperscissor.Result
	registerC   chan *client
	unregisterC chan *client
}

type client struct {
	id   string
	hub  *hub
	conn *websocket.Conn
	move rockpaperscissor.Move
}

func (c *client) Move() rockpaperscissor.Move { return c.move }
func (c *client) Name() string                { return c.id }
func (c *client) ResetMove()                  { c.move = 0 }

func (h *hub) run() {
	for {
		select {
		case c := <-h.registerC:
			h.players[c.id] = c
		case c := <-h.unregisterC:
			delete(h.players, c.id)
		case msg := <-h.moveC:
			c := h.players[msg.ID]
			c.move = msg.Move
			if ok := h.isPlayersReady(); !ok {
				continue
			}

			opponent, err := h.getOpponent(c.Name())
			if err != nil {
				log.Println(err)
			}

			res := rockpaperscissor.Play(c, opponent)
			h.history = append(h.history, res)
			h.resultsC <- res
		}
	}
}

func (h *hub) isPlayersReady() bool {
	for _, p := range h.players {
		if p.move == 0 {
			return false
		}
	}
	return true
}

func (h *hub) getOpponent(id string) (*client, error) {
	for _, p := range h.players {
		if p.Name() != id {
			return p, nil
		}
	}
	return nil, errors.New("couldn't find any opponents in the map.")
}

type ReadMessage struct {
	ID   string                `json:"id"`
	Move rockpaperscissor.Move `json:"move"`
}

func (c *client) read() {
	defer func() {
		c.hub.unregisterC <- c
		c.conn.Close()
	}()

	for {
		var req ReadMessage
		err := c.conn.ReadJSON(&req)
		if err != nil {
			log.Printf("Error, can't seem to read the message: %v", err)
			break
		}

		c.hub.moveC <- req
	}
}

func (c *client) write() {
	ticker := time.NewTicker(time.Second * 5)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case res := <-c.hub.resultsC:
			fmt.Println(res.Winner.Name())
		}

	}
}
