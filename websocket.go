package main

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	// currentStatus     nomachineStatus
	upgrader = websocket.Upgrader{}
	// statusBroadcaster = SetupBroadcaster[nomachineStatus]()
)

func serveWs(c echo.Context) (err error) {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		doExit := false

		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			c.Logger().Error(err)
			// doExit = true
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			_, ok := err.(*websocket.CloseError)
			if ok {
				c.Logger().Info("ws connection closed")
				doExit = true
			}
			c.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)

		if doExit {
			break
		}
	}
	return
}

// func writer(ws *websocket.Conn, subscriber *Subscriber[nomachineStatus]) {
// 	statusTicker := time.NewTicker(10 * time.Second)
// 	defer func() {
// 		statusTicker.Stop()
// 		ws.Close()
// 	}()
// 	for {
// 		select {
// 		case <-statusTicker.C:
// 			writeStatusToWebsocket(ws)
// 		}
// 	}
// }

// // writeStatusToWebsocket writes the status to one web socket
// func writeStatusToWebsocket(ws *websocket.Conn) {
// 	status := currentStatus.Load()
// 	statusJsonBytes, err := json.Marshal(status)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	ws.SetWriteDeadline(time.Now().Add(time.Minute))
// 	if err := ws.WriteMessage(websocket.TextMessage, statusJsonBytes); err != nil {
// 		log.Println(err)
// 	}
// }

// func checkStatusTickerHandler(ticker time.Ticker, statusUpdatesChannel chan nomachineStatus) {
// 	for {
// 		<-ticker.C

// 	}
// }
