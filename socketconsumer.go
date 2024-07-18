package socketconsumer

import (
	"log"
	"net/http"
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
