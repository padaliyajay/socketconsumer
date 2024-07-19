package socketconsumer

import (
	"log"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request, cons ConsumerInterface) {
	conn, err := hub.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	controller := NewConsumerController(hub, conn, &cons)

	if cons.Connect() {
		controller.open()
		go controller.WebsocketReceive()
		go controller.WebsocketSend()
	} else {
		conn.Close()
	}
}
