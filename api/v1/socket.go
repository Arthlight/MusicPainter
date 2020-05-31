package v1

import (
	"Spotify-Visualizer/models"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		HandshakeTimeout:  20,
		ReadBufferSize:    2048,
		WriteBufferSize:   2048,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func createNewWebsocket(w http.ResponseWriter, r *http.Request) (*models.WebSocket, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil{
		fmt.Printf("An error occured while trying to upgrade the connection: %v", err)
		return nil, err
	}

	ws := &models.WebSocket{
		Conn:  conn,
		Out:    make(chan []byte),
		In:     make(chan []byte),
		Events: make(map[string]models.EventHandler),
	}

}









// TODO: Eventuell eine goroutine schreiben die vom refresh token handler gecalled wird
// TODO: und diese goroutine called dann eine andere goroutine die die spotify api
// TODO: called und werte berechnet und die dann in derselbigen goroutine immer checkt ob der client noch
// TODO: connected ist, wenn ja, wird durch einen channel zurueck an die vorherige gorutine die berechneten
// TODO: werte gesendet, wenn nein, wird der channel geclosed und die vorherige goroutine stopt den for loop
// TODO: und exited auch. Die werte werden von der ersten gorutine innerhalb des for loop dann ueber die
// TODO: socket ans frontend gesendet.




