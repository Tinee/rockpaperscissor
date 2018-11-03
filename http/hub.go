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
	players map[string]*listener
	history []rockpaperscissor.Result

	moveC       chan ReadMessage
	resultsC    chan rockpaperscissor.Result
	registerC   chan *listener
	unregisterC chan *listener
}

type listener struct {
	id   string
	hub  *hub
	conn *websocket.Conn
	move rockpaperscissor.Move
}

func (c *listener) Move() rockpaperscissor.Move { return c.move }
func (c *listener) Name() string                { return c.id }
func (c *listener) ResetMove()                  { c.move = 0 }

func (h *hub) run() {

	for {
		select {
		case c := <-h.registerC:
			h.players[c.id] = c
			go c.read()
			go c.write()
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

func (h *hub) getOpponent(id string) (*listener, error) {
	for _, p := range h.players {
		if p.Name() != id {
			return p, nil
		}
	}
	return nil, errors.New("couldn't find any opponents in the map")
}

// ReadMessage is the message the peers will send us.
type ReadMessage struct {
	ID   string                `json:"id"`
	Move rockpaperscissor.Move `json:"move"`
}

func (c *listener) read() {
	defer func() {
		c.hub.unregisterC <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		return nil
	})

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

func (c *listener) write() {
	ticker := time.NewTicker(time.Second * 5)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case res := <-c.hub.resultsC:
			fmt.Println(res.Winner.Name())
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
