# Socket Consumer
Create a socket consumer which listen incoming socket requests and manage it. this package uses [gorilla websocket](https://github.com/gorilla/websocket) and add new functionality like creating consumer and user room over it. it is made by inspiration from python [django channels](https://github.com/django/channels) library


### Installation

    go get github.com/padaliyajay/socketconsumer

## Example

#### initialization

    upgrader := websocket.Upgrader{
    	ReadBufferSize:  1024,
    	WriteBufferSize: 1024,
    	CheckOrigin: func(r *http.Request) bool {
    		return true
    	},
    }
    
    hub := socketconsumer.NewHub(upgrader)

#### ChatConsumer.go

    type ChatConsumer struct {
    	groupName string
    	
    	*socketconsumer.Consumer
    }
    
    func ServeChatConsumer(hub *socketconsumer.Hub) func(w http.ResponseWriter, r *http.Request) {
    	return func(w http.ResponseWriter, r *http.Request) {
    		socketconsumer.ServeWS(hub, w, r, &ChatConsumer{
    			groupName: "my_group",
    			Consumer: socketconsumer.NewConsumer(hub),
    		})
    	}
    }
    
    // Aurhenticate user
    func (c *ChatConsumer) Connect() bool {
    	return true
    }
    
    // Create group
    func (c *ChatConsumer) Accept() {
    	c.GroupAdd(c.groupName)
    }
    
    // Close group
    func (c *ChatConsumer) Disconnect() {
    	c.GroupDiscard(c.groupName)
    }
    
    // Receive message
    func (c *ChatConsumer) Receive(message *socketconsumer.Message) {
    	if message.Type == "message" {
    		// Send message to all in group
    		c.GroupSend(c.groupName, message)
    
    		// Send message to others in group except me
    		c.GroupSendOthers(c.groupName, message)
    
    		// Send to me only
    		c.Send(message)
    
    		// Send to specific user
    		c.SendTo("consumer_id", message)
    	}
    }

#### Assign consumer to request

    http.HandleFunc("/ws", ServeChatConsumer(hub))
