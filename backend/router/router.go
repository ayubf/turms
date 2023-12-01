package router

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"log"
	"net/http"
	"time"
	"turmsapi/util"
)

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	util.WrapPostMethod(w, r, func(w http.ResponseWriter, r *http.Request) {
		var reqBody struct {
			UserName  string
			RoomName  string
			TimeLimit int
		}

		bod, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			log.Fatal(err)
		}
		err = json.Unmarshal(bod, &reqBody)
		if err != nil {
			http.Error(w, "Error decoding JSON", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")

		roomCode := util.CreateNewCode(6)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"created":  true,
			"code":     roomCode,
			"roomName": reqBody.RoomName,
			"username": reqBody.UserName,
		})

		ctx, client := util.InitClient()

		rooms := client.Database("turmsdb").Collection("rooms")

		rooms.InsertOne(ctx, bson.M{
			"roomCreatorName": reqBody.UserName,
			"roomName":        reqBody.RoomName,
			"code":            roomCode,
			"messages":        []util.Message{},
			"expirationAt":    time.Now().Local().Add(time.Duration(reqBody.TimeLimit) * time.Minute),
		})

		util.AuxMessage("/createRoom")

	})
}

func JoinRoom(w http.ResponseWriter, r *http.Request) {

	util.EnableCors(&w)

	ctx, client := util.InitClient()

	rooms := client.Database("turmsdb").Collection("rooms")
	code := r.URL.Query().Get("code")

	filter := bson.M{"code": code}

	var specificRoom util.Room
	err := rooms.FindOne(ctx, filter).Decode(&specificRoom)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println(err)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"roomExists": false,
			})
			util.AuxMessage("/joinRoom")
			return
		}
		panic(err)
	}

	expireCheck := func() {
		if specificRoom.ExpirationAt.Before(time.Now()) {
			w.WriteHeader(http.StatusForbidden)
			util.AuxMessage("/joinRoom")
			fmt.Println("Expired")
			_, err = rooms.DeleteOne(ctx, filter)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			util.AuxMessage("/joinRoom")
			fmt.Println("Active")
		}
	}

	expireCheck()

	util.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := util.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Could not open WebSocket connection")
		return
	}
	err = util.SendJSON(conn, specificRoom)
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
		var receivedMessage struct {
			UserName string
			Message  string
		}
		err = json.Unmarshal(p, &receivedMessage)
		if err != nil {
			log.Println("Error unmarshaling JSON:", err)
			return
		}
		newMessage := util.Message{
			Message: receivedMessage.Message,
			Author:  receivedMessage.UserName,
			Time:    time.Now().Format(time.UnixDate),
		}

		expireCheck()
		json.NewEncoder(w).Encode(map[string]interface{}{
			"alert":      "Room Closed",
			"roomExists": false,
		})

		update := bson.D{{Key: "$push", Value: bson.D{{Key: "messages", Value: newMessage}}}}
		_, err = rooms.UpdateOne(ctx, filter, update)
		if err != nil {
			fmt.Println("Error updating document:", err)
			return
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			return
		}

	}

}

func CheckRooms() {
	uptimeTicker := time.NewTicker(1 * time.Hour)
	ctx, client := util.InitClient()

	filter := bson.M{"$expr": bson.M{"$lt": []interface{}{"$expirationAt", time.Now()}}}
	rooms := client.Database("turmsdb").Collection("rooms")

	for range uptimeTicker.C {
		_, err := rooms.DeleteMany(ctx, filter)
		if err != nil {
			log.Fatal(err)
		}
	}

}
