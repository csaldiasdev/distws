package wshub

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
)

type WsConnection interface {
	SendMessage(msg []byte)
}

func NewConnection(wsConn *websocket.Conn, unRegConnChan chan<- uuid.UUID) (WsConnection, uuid.UUID) {

	connId := uuid.New()

	conn := &connection{
		wsConn:     wsConn,
		send:       make(chan []byte),
		unregister: unRegConnChan,
		id:         connId,
	}

	conn.wsConn.SetReadLimit(maxMessageSize)
	conn.wsConn.SetReadDeadline(time.Now().Add(pongWait))
	conn.wsConn.SetPongHandler(func(string) error { conn.wsConn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	go conn.readPump()
	go conn.writePump()

	return conn, connId
}

type connection struct {
	closeOnce  sync.Once
	pingTicker *time.Ticker
	wsConn     *websocket.Conn
	send       chan []byte
	unregister chan<- uuid.UUID
	id         uuid.UUID
}

func (c *connection) SendMessage(msg []byte) {
	c.send <- msg
}

func (c *connection) writePump() {
	c.pingTicker = time.NewTicker(pingPeriod)
	defer c.closeConnection()

	for {
		select {
		case message, ok := <-c.send:
			c.wsConn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.wsConn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.wsConn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-c.pingTicker.C:
			c.wsConn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.wsConn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *connection) readPump() {
	defer c.closeConnection()
	for {
		_, _, err := c.wsConn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error().Err(err).Msg("Unexpected websocket close error")
			}
			break
		}
	}
}

func (c *connection) closeConnection() {
	c.closeOnce.Do(func() {
		c.writeConnectionLog("closing websocket connection")
		c.unregister <- c.id
		close(c.send)
		c.wsConn.Close()
		if c.pingTicker != nil {
			c.pingTicker.Stop()
		}
	})
}

func (c *connection) writeConnectionLog(message string) {
	log.Info().
		Str("connectionId", c.id.String()).
		Msg(message)
}
