package main

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	currentStatus nomachineStatus
	upgrader      = websocket.Upgrader{}
	// statusBroadcaster = SetupBroadcaster[nomachineStatus]()
)

func serveWs(c echo.Context) (err error) {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			c.Logger().Error(err)
		}
		return
	}

	go func() {
		defer ws.Close()
		writer(ws, c)
	}()
	// reader(ws)

	return
}

func writer(ws *websocket.Conn, c echo.Context) {
	statusTicker := time.NewTicker(3 * time.Second)
	defer statusTicker.Stop()

	lastStatus, err := getStatus()
	if err != nil {
		c.Logger().Error(err)
		return
	}
	if err := writeStatusToWebsocket(lastStatus, ws, c); err != nil {
		return
	}

	for {
		select {
		case <-statusTicker.C:
			newStatus, err := getStatus()
			if err != nil {
				c.Logger().Error(err)
				break
			}
			if newStatus != lastStatus {
				lastStatus = newStatus
				if err := writeStatusToWebsocket(lastStatus, ws, c); err != nil {
					break
				}
			}
		}
	}
}

// writeStatusToWebsocket writes the status to one web socket
func writeStatusToWebsocket(status nomachineStatus, ws *websocket.Conn, c echo.Context) (err error) {
	statusJsonBytes, err := json.Marshal(status)
	if err != nil {
		c.Logger().Error(err)
		return
	}
	ws.SetWriteDeadline(time.Now().Add(time.Minute))
	if err := ws.WriteMessage(websocket.TextMessage, statusJsonBytes); err != nil {
		c.Logger().Error(err)
	}

	return
}

// func checkStatusTickerHandler(ticker time.Ticker, statusUpdatesChannel chan nomachineStatus) {
// 	for {
// 		<-ticker.C

// 	}
// }
