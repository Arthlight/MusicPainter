package models

import (
	"encoding/json"
)

type EventHandler func(event *Event)

type Event struct {
	Name string `json:"name"`
	Content interface{} `json:"content"`
}

func NewEventFromBinary(rawData []byte) (*Event,error) {
	event := new(Event)
	err := json.Unmarshal(rawData, event)
	return event, err
}

func (e *Event) ToBinary() ([]byte, error) {
	return json.Marshal(*e)
}
