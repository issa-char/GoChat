package internal

import (
  "log"
  "net/http"
  "github.com/gorilla/websocket"
}

// and upgrader to handle socket connnections
var upgrader = websocket.Upgrader {
  ReadBufferSize: 1024,
  WriteBufferSize: 1024,
  CheckOrigin: func(r *http.Request) bool {
    return true
  },
}

// handling incoming websocket connections
func handleConnections(w http.ResponseWriter, r *http.Request) {
  ws, err := upgrader.Upgrader(w, r, nil)
  if err != nil {
    log.Printf("unable to upgrade: %v", err)
    break
  }
  defer ws.Close()

  for {
    // read message from client
    _, message, err := ws.ReadMessage()
    if err != nil {
      log.Printf("error reading message: %v", err)
      break
    }
    log.Printf("Received: %v", message)

    // broadcast message back to client
    err = ws.WriteMessage(websocket.TextMessage, message)
    if err != nil {
      log.Printf("error writing message: %v", err)
      break
    }
  }
}

