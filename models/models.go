package models

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

type RefreshToken struct {
	Token string `json:"token"`
}

type Canvas struct {
	X int
	Y int
}

type EventHandler func(*Event)

type WebSocket struct {
	Conn *websocket.Conn
	Out chan []byte
	In chan []byte
	Events map[string]EventHandler
}

func (w *WebSocket) Reader() {
	defer func() {
		fmt.Printf("Error while trying to close Reader: %v", w.Conn.Close())
	}()

	for {
		_, message, err := w.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseGoingAway) {
				fmt.Printf("Error while trying to read message: %v", err)
			}
			break
		}

		event, err := newEventFromBinary(message)
		if err != nil {
			fmt.Printf("Error while trying to convert message from binary to struct: %v", err)
		} else {
			fmt.Printf("Received from Frontend: %v", event.Content)
		}

		if eventHandler, ok := w.Events[event.Name]; ok {
			eventHandler(event)
		}
	}
}

func (w *WebSocket) Writer() {

}

type Event struct {
	Name string `json:"name"`
	Content interface{} `json:"content"`
}

func newEventFromBinary(rawData []byte) (*Event,error) {
	event := new(Event)
	err := json.Unmarshal(rawData, event)
	return event, err
}

func (e *Event) toBinary() ([]byte, error) {
	return json.Marshal(e)
}
