package util

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
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
	RoomCreatorName string    `json:"roomCreatorName"`
	RoomName        string    `json:"roomName"`
	Code            string    `json:"code"`
	Messages        []Message `json:"messages"`
	ExpirationAt    time.Time `json:"expirationAt"`
}

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func SendJSON(conn *websocket.Conn, roomInfo Room) error {
	response, err := json.Marshal(roomInfo)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, response)
}

func InitClient() (context.Context, *mongo.Client) {
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

	return ctx, client

}

func AuxMessage(r string) {
	t := time.Now()
	fmt.Printf("-> %v:%v @ %v  \n", p, r, t.Format(time.UnixDate))
}

func EnableCors(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	if origin == "http://localhost:3000" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
}

func CreateNewCode(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func WrapMethod(method string, w http.ResponseWriter, r *http.Request, f func(http.ResponseWriter, *http.Request)) {
	EnableCors(w, r)
	if r.Method == method {
		f(w, r)
	} else {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("403 - Forbidden Method"))
		AuxMessage("/createSession")
	}
}

func CheckRooms() {
	uptimeTicker := time.NewTicker(5 * time.Second)
	ctx, client := InitClient()

	filter := bson.M{"$expr": bson.M{"$lt": []interface{}{"$expirationAt", time.Now()}}}
	rooms := client.Database("turmsdb").Collection("rooms")

	for range uptimeTicker.C {
		_, err := rooms.DeleteMany(ctx, filter)
		if err != nil {
			log.Fatal(err)
		}
	}

}
