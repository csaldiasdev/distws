package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/csaldiasdev/distws/internal/httpserver/jwt"
	"github.com/csaldiasdev/distws/internal/wshub"
	"github.com/go-chi/chi/v5"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type httpServer struct {
	Hub       *wshub.Hub
	jwtIssuer string
	jwtSecret string
}

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *httpServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Message string `json:"message"`
	}{"Hello world!"}

	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpServer) handleMessageToUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	userId, err := uuid.Parse(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body := struct {
		Message string `json:"message"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.Hub.MessageToUser(userId, []byte(body.Message))
	w.WriteHeader(http.StatusOK)
}

func (s *httpServer) handleWs(w http.ResponseWriter, r *http.Request) {

	keys, ok := r.URL.Query()["token"]

	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenString := keys[0]

	payload, err := jwt.ValidateToken(tokenString, s.jwtIssuer, s.jwtSecret)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Error().Err(err).Msg("Error on upgrade connection to the WebSocket protocol")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.Hub.AddWebSocketConnection(conn, payload.Sub)
}

func NewHTTPServer(hub *wshub.Hub) *http.Server {
	httpsvr := httpServer{
		Hub:       hub,
		jwtIssuer: "http://distributedws",
		jwtSecret: "distributedws",
	}

	chiRouter := chi.NewRouter()

	chiRouter.Get("/", httpsvr.handleRoot)
	chiRouter.Post("/api/user/{id}/message", httpsvr.handleMessageToUser)
	chiRouter.Get("/ws", httpsvr.handleWs)

	return &http.Server{
		Handler: chiRouter,
	}
}
