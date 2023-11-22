package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	p = 8080
)

type InputMessage struct {
	Author  string `json:"author"`
	Message string `json:"message"`
}

type Message struct {
	Author  string
	Message string
	Time    string
}

type Room struct {
	RoomCreatorName string
	RoomName        string
	Code            string
	Messages        []Message
}

func createNewCode(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func auxMessage(r string) string {
	t := time.Now()
	return fmt.Sprintf("-> %v:%v @ %v", p, r, t.Format(time.UnixDate))
}

// Creates Room From Homepage
func createRoomRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Println(auxMessage("/createRoom"))
	if r.Method == http.MethodPost {
		bod, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}
		var reqBody struct {
			UserName string
			RoomName string
		}

		fmt.Println(string(bod))
		err = json.Unmarshal(bod, &reqBody)
		if err != nil {
			http.Error(w, "Error decoding JSON", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")

		roomCode := createNewCode(6)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"created":  true,
			"code":     roomCode,
			"roomName": reqBody.RoomName,
			"username": reqBody.UserName,
		})

		ctx := context.TODO()

		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatal(err)
		}
		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}

		rooms := client.Database("turmsdb").Collection("rooms")

		rooms.InsertOne(ctx, &Room{
			RoomCreatorName: reqBody.UserName,
			RoomName:        reqBody.RoomName,
			Code:            roomCode,
			Messages:        []Message{},
		})
		fmt.Println(auxMessage("/"))

	} else {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("403 - Forbidden Method"))
		auxMessage("/")
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Join Rooms
func joinRoom(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	rooms := client.Database("turmsdb").Collection("rooms")
	code := r.URL.Query().Get("code")

	filter := bson.D{primitive.E{Key: "code", Value: code}}

	var specificRoom Room
	err = rooms.FindOne(ctx, filter).Decode(&specificRoom)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"roomExists": false,
			})
			auxMessage("/")
			return
		}
		panic(err)
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open WebSocket connection", http.StatusBadRequest)
		return
	}
	err = sendJSON(conn, specificRoom)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		var receivedMessage Message
		err = json.Unmarshal(p, &receivedMessage)
		if err != nil {
			log.Println("Error unmarshaling JSON:", err)
			return
		}
		newMessage := Message{
			Message: receivedMessage.Message,
			Author:  receivedMessage.Author,
			Time:    time.Now().Format(time.UnixDate),
		}
		update := bson.D{{Key: "$push", Value: bson.D{{Key: "messages", Value: newMessage}}}}
		_, err = rooms.UpdateOne(context.Background(), filter, update)
		if err != nil {
			fmt.Println("Error updating document:", err)
			return
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			return
		}

	}

}

func sendJSON(conn *websocket.Conn, roomInfo Room) error {
	response, err := json.Marshal(roomInfo)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, response)
}

func main() {
	t := time.Now()
	http.HandleFunc("/createroom", createRoomRoute)
	http.HandleFunc("/joinroom", joinRoom)

	fmt.Printf("-> Server Running on Port: %v\n-> %v\n", p, t.Format(time.UnixDate))
	err := http.ListenAndServe(":"+fmt.Sprint(p), nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server Closed")
	} else if err != nil {
		fmt.Printf("Error while starting server: %v\n", err)
		os.Exit(1)
	}

}
