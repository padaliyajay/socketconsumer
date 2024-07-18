package socketconsumer

type ConsumerInterface interface {
	Connect() bool
	Accept()
	Disconnect()
	Receive(message *Message)
	Send(message *Message)
	GetId() string
}

type Consumer struct {
	// Unique id of the consumer
	Id string

	hub *Hub
}

func NewConsumer(hub *Hub) *Consumer {
	return &Consumer{
		Id:  RandomString(10),
		hub: hub,
	}
}

func (c *Consumer) GetId() string {
	return c.Id
}

// Connect function for accept or reject the connection. used to authenticate the connection
func (c *Consumer) Connect() bool {
	return true
}

// Disconnect function for disconnect task
func (c *Consumer) Disconnect() {}

// Accept function for accept the connection
func (c *Consumer) Accept() {}

// receive message
func (c *Consumer) Receive(message *Message) {}

// Send message
func (c *Consumer) Send(message *Message) {
	c.hub.Send(c.Id, message)
}

// Send message to other consumer
func (c *Consumer) SendTo(id string, message *Message) {
	c.hub.Send(id, message)
}

// Add group
func (c *Consumer) GroupAdd(name string) {
	c.hub.GroupAdd(name, c)
}

// Discard group
func (c *Consumer) GroupDiscard(name string) {
	c.hub.GroupDiscard(name, c)
}

// Send to group
func (c *Consumer) GroupSend(name string, message *Message) {
	c.hub.GroupSend(name, message)
}

// Send to group except self
func (c *Consumer) GroupSendOthers(name string, message *Message) {
	c.hub.GroupSendExcept(name, message, c)
}
