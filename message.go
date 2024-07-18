package socketconsumer

import (
	"encoding/json"
	"log"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Convert message to json
func (m *Message) ToJson() []byte {
	jsonData, err := json.Marshal(m)

	if err != nil {
		log.Println(err)
		return []byte{}
	}

	return jsonData
}

// Create new message
func NewMessage(t string, d interface{}) *Message {
	return &Message{
		Type: t,
		Data: d,
	}
}

// from json
func NewMessageFromJson(data []byte) *Message {
	m := &Message{}
	err := json.Unmarshal(data, m)

	if err != nil {
		log.Println(err)
		return nil
	}

	return m
}
