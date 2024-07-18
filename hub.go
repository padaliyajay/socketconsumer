package socketconsumer

import (
	"github.com/gorilla/websocket"
)

type Hub struct {
	// Registered consumers
	consumers map[string]*ConsumerController

	// Registered groups
	groups map[string]*group

	// upgrader
	upgrader websocket.Upgrader
}

// newHub creates a new hub
func NewHub(upgrader websocket.Upgrader) *Hub {
	return &Hub{
		consumers: make(map[string]*ConsumerController),
		groups:    make(map[string]*group),
		upgrader:  upgrader,
	}
}

// Register consumer
func (hub *Hub) RegisterConsumer(c *ConsumerController) {
	hub.consumers[c.getId()] = c
}

// Unregister consumer
func (hub *Hub) UnregisterConsumer(c *ConsumerController) {
	delete(hub.consumers, c.getId())
}

// Get consumer by id
func (hub *Hub) GetConsumer(id string) *ConsumerController {
	if c, ok := hub.consumers[id]; ok {
		return c
	}
	return nil
}

// Has consumer
func (hub *Hub) HasConsumer(id string) bool {
	_, ok := hub.consumers[id]
	return ok
}

// Send to consumer
func (hub *Hub) Send(id string, message *Message) {
	if c, ok := hub.consumers[id]; ok {
		c.send(message)
	}
}

// Add consumer to group
func (hub *Hub) GroupAdd(id string, c *Consumer) {
	if _, ok := hub.groups[id]; !ok {
		hub.groups[id] = newGroup(id)
	}
	hub.groups[id].addConsumer(c)
}

// Discard consumer from group and delete empty group
func (hub *Hub) GroupDiscard(id string, c *Consumer) {
	hub.groups[id].removeConsumer(c)
	if hub.groups[id].isEmpty() {
		delete(hub.groups, id)
	}
}

// Send to group
func (hub *Hub) GroupSend(id string, message *Message) {
	hub.groups[id].broadcast(message)
}

// Send to group except consumer
func (hub *Hub) GroupSendExcept(id string, message *Message, except *Consumer) {
	hub.groups[id].broadcastExcept(message, except)
}
