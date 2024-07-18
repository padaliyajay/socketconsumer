package socketconsumer

type group struct {
	// Unique id of the group
	id string

	// Registered consumers
	consumers map[*Consumer]bool
}

// newGroup creates a new group
func newGroup(id string) *group {
	return &group{
		id:        id,
		consumers: make(map[*Consumer]bool),
	}
}

// addConsumer adds a consumer to the group
func (g *group) addConsumer(c *Consumer) {
	g.consumers[c] = true
}

// removeConsumer removes a consumer from the group
func (g *group) removeConsumer(c *Consumer) {
	delete(g.consumers, c)
}

// Check group is empty
func (g *group) isEmpty() bool {
	return len(g.consumers) == 0
}

// broadcast sends a message to all consumers in the group
func (g *group) broadcast(message *Message) {
	for c := range g.consumers {
		c.Send(message)
	}
}

// broadcastExcept sends a message to all consumers in the group except the specified consumer
func (g *group) broadcastExcept(message *Message, except *Consumer) {
	for c := range g.consumers {
		if c != except {
			c.Send(message)
		}
	}
}
