package models

import (
	"encoding/json"
	"fmt"
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

func (e *Event) ToBinary() []byte {
	event, err := json.Marshal(*e)
	if err != nil {
		fmt.Println("Error while trying to marshal event: ", err)
	}

	return event
}
