package http

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// GameHandler represents the core handler
type GameHandler struct {
	hubs map[string]*hub

	*mux.Router
	sync.RWMutex
	websocket.Upgrader
}

// NewGameHandler sets up
func NewGameHandler() *GameHandler {
	mux := mux.NewRouter()
	gh := &GameHandler{
		Router: mux,
		hubs:   make(map[string]*hub),
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}

	mux.HandleFunc("/game", gh.AddGame).Methods("POST")
	mux.HandleFunc("/ws", gh.JoinGame).Methods("GET")

	return gh
}

// NewGameRequest is the request body when adding a new game.
type NewGameRequest struct {
	PlayerID string `json:"id"`
	GameName string `json:"gameName"`
}

// NewGameResponse is the request body when adding a new game.
type NewGameResponse struct {
	PlayerID string `json:"id"`
	GameName string `json:"gameName"`
}

// AddGame takes a request and adds it.
func (g *GameHandler) AddGame(w http.ResponseWriter, r *http.Request) {
	conn, err := g.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "could not upgrade the request", http.StatusBadRequest)
		return
	}

	var req NewGameRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "failed to parse body", http.StatusBadRequest)
		return
	}
	hub := hub{
		players:     make(map[string]*client),
		moveC:       make(chan ReadMessage),
		registerC:   make(chan *client),
		unregisterC: make(chan *client),
	}

	g.RLock()
	defer g.RUnlock()
	g.hubs[req.GameName] = &hub

	hub.registerC <- &client{
		conn: conn,
		id:   req.PlayerID,
	}

	bs, err := json.Marshal(&NewGameResponse{
		GameName: req.GameName,
		PlayerID: req.PlayerID,
	})
	if err != nil {
		http.Error(w, "failed to parse response body", http.StatusBadRequest)
		return
	}

	go hub.run()

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}

// JoinGameRequest is the request body when joining a game.
type JoinGameRequest struct {
	PlayerID string `json:"id"`
	GameName string `json:"gameName"`
}

func (g *GameHandler) JoinGame(w http.ResponseWriter, r *http.Request) {

	conn, err := g.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "could not upgrade the request", http.StatusBadRequest)
		return
	}

	var req JoinGameRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "could not parse the incoming request", http.StatusBadRequest)
		return
	}

	g.RLock()
	defer g.RUnlock()
	hub, ok := g.hubs[req.GameName]
	if !ok {
		http.Error(w, "no game with that name was found.", http.StatusNotFound)
	}

	hub.registerC <- &client{
		conn: conn,
		hub:  hub,
		id:   req.GameName,
	}

	w.WriteHeader(http.StatusOK)
}
