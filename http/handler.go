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

// NewGameHandler sets up a Http.Handler
func NewGameHandler() *GameHandler {
	mux := mux.NewRouter()
	gh := &GameHandler{
		Router: mux,
		hubs:   make(map[string]*hub),
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// Fix this so it's more secure.
				return true
			},
		},
	}

	mux.HandleFunc("/game", gh.AddGame).Methods("POST")
	mux.HandleFunc("/ws/{id}/{game}", gh.JoinGame).Methods("GET")

	return gh
}

// NewGameRequest is the request body when adding a new game.
type NewGameRequest struct {
	GameName string `json:"gameName"`
}

// NewGameResponse is the response body when adding a new game.
type NewGameResponse struct {
	GameName string `json:"gameName"`
}

// AddGame tries to create a new game,
// it also put's you in as a player and sets up a websocket connection.
func (g *GameHandler) AddGame(w http.ResponseWriter, r *http.Request) {
	var req NewGameRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "failed to parse body", http.StatusBadRequest)
		return
	}
	hub := hub{
		players:     make(map[string]*listener),
		moveC:       make(chan ReadMessage),
		registerC:   make(chan *listener),
		unregisterC: make(chan *listener),
	}

	g.RLock()
	defer g.RUnlock()
	g.hubs[req.GameName] = &hub

	bs, err := json.Marshal(&NewGameResponse{
		GameName: req.GameName,
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

// JoinGame upgrades you to a websocket connection,
// and tries to find the game you asked for and put you in as a player there.
func (g *GameHandler) JoinGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	game := vars["game"]

	g.RLock()
	defer g.RUnlock()
	hub, ok := g.hubs[game]
	if !ok {
		http.Error(w, "no game with that name was found.", http.StatusNotFound)
		return
	}

	conn, err := g.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hub.registerC <- &listener{
		conn: conn,
		hub:  hub,
		id:   id,
	}

	w.WriteHeader(http.StatusOK)
}
