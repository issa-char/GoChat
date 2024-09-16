package internal

import (
  "log"
  "net/http"
  "time"
  "github.com/golang-jwt/v5"
  "github.com/gorilla/websocket"
)

var jwtSecret = []byte("your_jwt_secret_key")

// generate a JWT token for a given username
func generateJWT(username string) (string, error) {
  claims := jwt.MapClaims{
    "username": username,
    "exp": time.Now().Add(time.Hour * 72).Unix(),
  }
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  return token.SignedString(jwtSecret)
}

// middleware to verify the JWT token in the Websocket request
func validateJWT(next http.HandleFunc) http.HandleFunc {
  return func(w http.ResponseWriter, r http.Request) {
    tokenString := r.URL.Query().Get("token")
    if tokenString == "" {
      http.Error(w, "missing token", http.StatusUnauthorized)
      return
    }

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
      if _, ok := token.Method.(*jwt.SigningMethodHMAC): !ok {
        return nil, http.ErrAbortHandler
      }
      return jwtSecret, nil
    })

    if err != nil || !token.valid {
      http.Error(w, "invalid token", http.StatusUnauthorized)
      return
    }

    // token is valid, continue to websocket connection
    next.ServeHTTP(w, r)
  }
}

