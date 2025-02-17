package socketconsumer

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type ConsumerController struct {
	cons *ConsumerInterface

	hub *Hub

	conn *websocket.Conn

	sendChan chan []byte
}

func NewConsumerController(hub *Hub, conn *websocket.Conn, cons *ConsumerInterface) *ConsumerController {
	return &ConsumerController{
		cons:     cons,
		hub:      hub,
		conn:     conn,
		sendChan: make(chan []byte),
	}
}

// open opens the consumer
func (cc *ConsumerController) open() {
	cc.hub.RegisterConsumer(cc)
	(*cc.cons).Accept()
}

// close closes the consumer
func (cc *ConsumerController) close() {
	(*cc.cons).Disconnect()
	cc.hub.UnregisterConsumer(cc)
	cc.conn.Close()
	close(cc.sendChan)
}

// Get id
func (cc *ConsumerController) getId() string {
	return (*cc.cons).GetId()
}

// Send message
func (cc *ConsumerController) send(message *Message) {
	cc.sendChan <- message.ToJson()
}

// read message from websocket
func (cc *ConsumerController) WebsocketReceive() {
	defer cc.close()

	cc.conn.SetReadDeadline(time.Now().Add(pongWait))
	cc.conn.SetPongHandler(func(string) error { cc.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := cc.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		(*cc.cons).Receive(NewMessageFromJson(message))
	}
}

// Write message to websoket
func (cc *ConsumerController) WebsocketSend() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		cc.conn.Close()
	}()

	for {
		select {
		case message, ok := <-cc.sendChan:
			if !ok {
				cc.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := cc.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Println(err)
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(cc.sendChan)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-cc.sendChan)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := cc.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println(err)
				return
			}
		}
	}
}
