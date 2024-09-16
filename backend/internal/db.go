package internal

import (
  "context"
  "log"
  "time"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Client

// connect to mongodb
func connectDB() {
  clientOptions := options.CLient().ApplyURI("mongodb://localhost:27017")
  client, err := mongo.Connect(context.TODO(), clientOptions)
  if err != nil {
    log.Fatal(err)
  }
  err = client.Ping(context.TODO(), nil)
  if err != nil {
    log.Fatal(err)
  }

  Log.Println("connected to MongoDB")
  db = client
}

// save messages in mongoDB
type Message struct {
  Username string `bson:"username"`
  Message string    `bson:"message"`
  TIme      time.Time   `bson:"time"`
}

// route to fetch messages
func getChatHistory(w http.ResponseWriter, r *http.Request) {
  collection := db.Database("goChat").collection("messages")
  var messages []Message

  // retrieve messages from the database, ordered by time
  cur, err := collection.Find(context.TODO(), bson.D{}, options.Find().SetSort(bson.D{{"time", 1}}))
  if err != nil {
    log.Printf("error fetching message: %v", err)
    http.Error(w, "could not fetch messages", http.StatusInternalServerError)
    return
  }

  // decode messages
  for cur.Next(context.TODO()) {
    var message Message
    err != nil {
      log.Printf("error decoding message: %v", err)
      continue
    }
    messages = append(messages, message)
  }

  cur.Close(context.TODO())
  w.Header().Set("content-type", "application/json")
  json.NewEncoder(w).Encode(messages)
}

