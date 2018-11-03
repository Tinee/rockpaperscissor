package http

import (
	"log"
	"rockpaperscissor/rockpaperscissor"
	"time"

	"github.com/gorilla/websocket"
)

type hub struct {
	listeners []*listener

	move       chan ReadMessage
	register   chan *listener
	unregister chan *listener
	ready      chan *listener
}

type listener struct {
	id    string
	ready bool
	move  rockpaperscissor.Move
	hub   *hub
	conn  *websocket.Conn

	outcome chan rockpaperscissor.Outcome
}

func (c *listener) ID() string                          { return c.id }
func (c *listener) Move() rockpaperscissor.Move         { return c.move }
func (c *listener) Decide(out rockpaperscissor.Outcome) { c.outcome <- out }

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.listeners = append(h.listeners, c)
			go c.read()
			go c.write()
		// case c := <-h.unregister:
		// delete(h.players, c.id)
		case msg := <-h.move:
			var l *listener
			for _, lsn := range h.listeners {
				if lsn.ID() == msg.ID {
					l = lsn
				}
			}
			l.move = msg.Move
			l.ready = true

			h.ready <- l
		case <-h.ready:
			if ready := h.isPartyReady(); !ready {
				// if the players isn't ready then do nothing.
				continue
			}
			first := h.listeners[0]
			second := h.listeners[1]
			rockpaperscissor.Play(first, second)
			first.ready = false
			second.ready = false
		}
	}
}

func (h *hub) isPartyReady() bool {
	if len(h.listeners) < 2 {
		// this is not a party?
		return false
	}

	for _, l := range h.listeners {
		if !l.ready {
			return false
		}
	}
	return true
}

// ReadMessage is the message the peers will send us.
type ReadMessage struct {
	ID   string                `json:"id"`
	Move rockpaperscissor.Move `json:"move"`
}

func (c *listener) read() {
	defer func() {
		c.hub.unregister <- c
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

		c.hub.move <- req
	}
}

type ResultResponse struct {
	Outcome rockpaperscissor.Outcome `json:"outcome"`
}

func (c *listener) write() {
	ticker := time.NewTicker(time.Second * 5)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case o := <-c.outcome:
			res := ResultResponse{Outcome: o}
			if err := c.conn.WriteJSON(&res); err != nil {
				log.Printf("Error, can't seem to write the message: %v", err)
				break
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
