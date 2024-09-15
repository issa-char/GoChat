package internal

import (
  "log"
  "net/http"
  "github.com/gorilla/websocket"
)


// websocket upgrader 
var upgrader = websocket.Upgrader {
  ReadBufferSize: 1024,
  WriteBufferSize: 1024,
  CheckOrigin: func(r *http.Request) bool {
    return true
  },
}

// a client struct
type Client struct {
  conn *websocket.Conn
  send chan []byte
}

// clientManager to handle all active clients
type ClientManager struct {
  clients map[*Client]bool
  broadcast chan []byte
  register chan *Client
  unregister chan *Client
}

var manager = ClientManager{
  clients: make(map[*Client]bool),
  broadcast: make(chan []byte),
  register: make(chan *Client),
  unregister: make(chan *Client),
}

// start the client manager to handle register, unregister, and broadcast events
func (manager *ClientManager) start() {
  for {
    select {
    case client := <-manager.register:
      manager.clients[client] = true
      log.Println("client registered")

    case client := <-manager.unregister:
      if _, ok := manager.clients[client]; ok {
        delete(manager.clients, client)
        close(client.send)
        log.Println("client unregistered")
      }
    case message := <-manager.broadcast:
      for client := range manager.clients {
        select {
        case client.send <- message:
        default:
          close(client.send)
          delete(manager.clients, client)
        }
      }
    }
  }
}


// handling websocket connections
func HandleConnections(w http.ResponseWriter, r *http.Request) {
  ws, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Printf("unable to upgrade: %v", err)
   return 
  }
  client := &Client{conn: ws, send:make(chan []byte)}

  manager.register <- client
  defer func() {
    manager.unregister <- client
    ws.Close()
  }()

  go client.writeMessages()
  client.readMessages()
}

// read messages from a websocket client
func (client *Client) readMessages() {
  for {
    _, message, err := client.conn.ReadMessage()
    if err != nil {
      log.Printf("error reading message: %v", err)
      manager.unregister <- client
      break
    }
    manager.broadcast <- message
  }
}

//write messages to a websocket client
func (client *Client) writeMessages() {
  for message := range client.send {
    err := client.conn.WriteMessage(websocket.TextMessage, message)
    if err != nil {
      log.Printf("error writing message: %v", err)
      manager.unregister <- client
      break
    }
  }
}


/*
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
  */

