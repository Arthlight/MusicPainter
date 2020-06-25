package models

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:    2048,
		WriteBufferSize:   2048,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type WebSocket struct {
	Conn *websocket.Conn
	Out chan []byte
	In <-chan []byte
	Events map[string]EventHandler
}

func CreateNewWebsocket(w http.ResponseWriter, r *http.Request) (*WebSocket, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil{
		fmt.Printf("An error occured while trying to upgrade the connection: %v", err)
		return nil, err
	}

	ws := &WebSocket{
		Conn:  conn,
		Out:    make(chan []byte),
		In:     make(<-chan []byte),
		Events: make(map[string]EventHandler),
	}

	go ws.Reader()
	go ws.Writer()

	return ws, nil

}

func (w *WebSocket) Reader() {
	defer func() {
		fmt.Println(w.Conn.Close())
	}()

	for {
		_, message, err := w.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseGoingAway) {
				fmt.Printf("Error while trying to read message: %v\n", err)
			}
			break
		}

		event, err := NewEventFromBinary(message)
		if err != nil {
			fmt.Printf("Error while trying to convert message from binary to struct: %v\n", err)
		} else {
			fmt.Printf("Received from Frontend: %v\n", event.Content)
		}

		if eventHandler, ok := w.Events[event.Name]; ok {
			eventHandler(event)
		}
	}
}

func (w *WebSocket) Writer() {
	for {
		select {
		case message, ok := <- w.Out:
			if !ok {
				err := w.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					panic(err)
				}
				return
			}
			writer, err := w.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			_, err = writer.Write(message)
			if err != nil {
				fmt.Println("Error while trying to send msg: ", err)
				panic(err)
			}
			err = writer.Close()
			if err != nil {
				fmt.Println("Error while trying to close connection: ", err)
				panic(err)
			}
		}
	}
}

func (w *WebSocket) On(eventName string, action EventHandler) *WebSocket {
	w.Events[eventName] = action
	return w
}