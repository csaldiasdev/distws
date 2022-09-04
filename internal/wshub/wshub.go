package wshub

import (
	"fmt"
	"sync"

	"github.com/csaldiasdev/distws/internal/repository"
	"github.com/csaldiasdev/distws/internal/wshub/rpc/client"
	"github.com/csaldiasdev/distws/internal/wshub/rpc/server"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type Hub struct {
	nodeId                  string
	repository              repository.Repository
	mu                      sync.Mutex
	rpcNodes                map[string]*client.HubRpcClient
	connections             map[uuid.UUID]WsConnection
	messageToConnectionChan chan messageToConnection
	registerChan            chan registerConnection
	unRegisterChan          chan uuid.UUID
}

type messageToConnection struct {
	id      uuid.UUID
	byteMsg []byte
}

type registerConnection struct {
	userId       uuid.UUID
	connectionId uuid.UUID
	conn         WsConnection
}

func NewHub(id string, rpcPort uint, repo repository.Repository) *Hub {
	h := &Hub{
		nodeId:                  id,
		repository:              repo,
		rpcNodes:                make(map[string]*client.HubRpcClient),
		connections:             make(map[uuid.UUID]WsConnection),
		messageToConnectionChan: make(chan messageToConnection),
		registerChan:            make(chan registerConnection),
		unRegisterChan:          make(chan uuid.UUID),
	}

	h.run()

	server.ListenAndServeHubRpc(rpcPort, func(userId string, message []byte) error {

		refConnections, err := h.repository.GetByUserId(uuid.MustParse(userId))

		if err != nil {
			return err
		}

		for _, v := range refConnections {
			if v.NodeId != h.nodeId {
				continue
			}

			h.messageToConnectionChan <- messageToConnection{uuid.MustParse(v.Id), message}
		}

		return nil
	})

	return h
}

func (h *Hub) AddHubNode(nodeId string, ip string, rpcPort uint) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	rpcCli, err := client.NewHubRpcClient(fmt.Sprintf("%s:%d", ip, rpcPort))

	if err != nil {
		return err
	}

	h.rpcNodes[nodeId] = rpcCli
	return nil
}

func (h *Hub) AddWebSocketConnection(wsConn *websocket.Conn) {
	conn, id := NewConnection(wsConn, h.unRegisterChan)
	h.registerChan <- registerConnection{conn: conn, connectionId: id}
}

func (h *Hub) MessageToUser(userId uuid.UUID, byteMsg []byte) error {
	refConnections, err := h.repository.GetByUserId(userId)

	if err != nil {
		return err
	}

	for _, v := range refConnections {
		if v.NodeId == h.nodeId {
			h.messageToConnectionChan <- messageToConnection{uuid.MustParse(v.Id), byteMsg}
		} else {
			rpcCli, ok := h.rpcNodes[v.NodeId]
			if ok {
				rpcCli.MessageToUser(userId.String(), byteMsg)
			}
		}
	}

	return nil
}

func (h *Hub) run() {
	go func() {
		for {
			select {
			case c := <-h.registerChan:
				h.connections[c.connectionId] = c.conn
				h.repository.Insert(c.connectionId, c.userId, uuid.MustParse(h.nodeId))

				log.Info().
					Str("userId", c.userId.String()).
					Str("connectionId", c.connectionId.String()).
					Str("nodeId", h.nodeId).
					Msg("User registered")

			case c := <-h.unRegisterChan:
				delete(h.connections, c)
				h.repository.DeleteConnection(c)

				log.Info().
					Str("connectionId", c.String()).
					Str("nodeId", h.nodeId).
					Msg("Connection unregistered")

			case m := <-h.messageToConnectionChan:
				wsc, found := h.connections[m.id]
				if found {
					wsc.SendMessage(m.byteMsg)
				}
			}
		}
	}()
}
