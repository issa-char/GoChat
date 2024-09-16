package main

import (
  "log"
  "net/http"
  "github.com/issa-char/goChat/backend/internal"
)


func main() {
  // connect to MongoDB
  connectDB()

  //server listens on port 3000 and handles websocket connections at /ws
  http.HandleFunc("/ws", internal.validateJWT(internal.HandleConnections))

  //
  go internal.manage.start()

  // route to get chat history
  http.HandleFunc("/history", getChatHistory)

  log.Println("starting websocket server on :3000")
  err := http.ListenAndServe(":3000", nil)
  if err != nil {
    log.Fatalf("server failed: %v", err)
  }
}

