package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// Client to talk to the websocket game server.
type Client struct {
	id       string
	wsDialer websocket.Dialer
	http     http.Client
	addr     string
	wsConn   *websocket.Conn
}

// NewClient sets up an *Client that is ready to use.
func NewClient(addr string, id string) *Client {
	c := &Client{
		http: http.Client{
			Timeout: 10 * time.Second,
		},
		id:   id,
		addr: addr,
		wsDialer: websocket.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: 45 * time.Second,
		},
	}
	return c
}

// Open opens a websocket connection on the given address.
func (c *Client) Open(game string) error {
	path := fmt.Sprintf("/ws/%s/%s", c.id, game)
	u := url.URL{Scheme: "ws", Host: c.addr, Path: path}

	conn, _, err := c.wsDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	c.wsConn = conn

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		fmt.Println(string(p))
	}
	return nil
}

// CreateGame makes a request and tries to create a new game.
func (c *Client) CreateGame(name string) (NewGameResponse, error) {
	var result NewGameResponse
	bs, err := json.Marshal(&NewGameRequest{GameName: name})
	if err != nil {
		return NewGameResponse{}, err
	}
	res, err := c.http.Post(fmt.Sprintf("http://%s/game", c.addr), "application/json", bytes.NewReader(bs))
	if err != nil {
		return NewGameResponse{}, err
	}

	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return NewGameResponse{}, err
	}

	return result, nil
}

// Close closes the underlaying websocket connection.
func (c *Client) Close() error {
	return c.wsConn.Close()
}
