package main

import (
  "log"
  "net/http"
  "github.com/issa-char/goChat/backend/internal"
)


func main() {
  //server listens on port 3000 and handles websocket connections at /ws
  http.HandleFunc("/ws", internal.HandleConnections)

  log.Println("starting websocket server on :3000")
  err := http.ListenAndServe(":3000", nil)
  if err != nil {
    log.Fatalf("server failed: %v", err)
  }
}

