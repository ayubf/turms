package router

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"turmsapi/util"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const SecretKey = "EXAMPLESECRETKEY"

func CreateSession(w http.ResponseWriter, r *http.Request) {
	util.WrapMethod("POST", w, r, func(w http.ResponseWriter, r *http.Request) {
		var reqBody struct {
			UserName string
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

		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Issuer:    reqBody.UserName,
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		})

		token, err := claims.SignedString([]byte(SecretKey))

		if err != nil {
			http.Error(w, "Error creating session", http.StatusInternalServerError)
		}

		c := http.Cookie{
			Name:     "token",
			Value:    token,
			Domain:   "localhost",
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
			Secure:   false,
			HttpOnly: true,
		}

		http.SetCookie(w, &c)
	})
}

func GetSession(w http.ResponseWriter, r *http.Request) {
	util.WrapMethod("GET", w, r, func(w http.ResponseWriter, r *http.Request) {

		tokenCookie, err := r.Cookie("token")
		if err == http.ErrNoCookie {
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"authenticated": "false",
			})
			return
		} else if err != nil {
			http.Error(w, "Error retrieving cookie", http.StatusInternalServerError)
		}

		token, err := jwt.ParseWithClaims(tokenCookie.Value, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

		if err != nil {
			http.Error(w, "Error reading cookie", http.StatusInternalServerError)
		}

		json.NewEncoder(w).Encode(token.Claims)

	})
}

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	util.WrapMethod("POST", w, r, func(w http.ResponseWriter, r *http.Request) {
		var reqBody struct {
			RoomName  string
			TimeLimit int
		}

		tokenCookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Error retrieving token cookie", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(tokenCookie.Value, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					http.Error(w, "Invalid token format", http.StatusUnauthorized)
					return
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					http.Error(w, "Token is expired or not yet valid", http.StatusUnauthorized)
					return
				} else {
					http.Error(w, "Error validating token", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "Error parsing token", http.StatusInternalServerError)
				return
			}
		}

		var username string

		fmt.Println(token.Claims)

		if tokenClaims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
			username = tokenClaims.Issuer
		} else {
			http.Error(w, "Error extracting token claims", http.StatusInternalServerError)
			return
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
			"username": username,
		})

		ctx, client := util.InitClient()

		rooms := client.Database("turmsdb").Collection("rooms")

		rooms.InsertOne(ctx, bson.M{
			"roomCreatorName": username,
			"roomName":        reqBody.RoomName,
			"code":            roomCode,
			"messages":        []util.Message{},
			"expirationAt":    time.Now().Local().Add(time.Duration(reqBody.TimeLimit) * time.Minute),
		})

		util.AuxMessage("/createRoom")

	})
}

func JoinRoom(w http.ResponseWriter, r *http.Request) {

	util.EnableCors(w, r)

	ctx, client := util.InitClient()

	rooms := client.Database("turmsdb").Collection("rooms")
	code := r.URL.Query().Get("code")

	filter := bson.M{"room": code}

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

	tokenCookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Error retrieving cookie", http.StatusInternalServerError)
	}

	token, err := jwt.ParseWithClaims(tokenCookie.Value, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		http.Error(w, "Error reading cookie", http.StatusInternalServerError)
	}

	tokenClaims := token.Claims.(jwt.MapClaims)

	username := tokenClaims["username"].(string)

	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		var receivedMessage struct {
			Message string
		}
		err = json.Unmarshal(p, &receivedMessage)
		if err != nil {
			log.Println("Error unmarshaling JSON:", err)
			return
		}
		newMessage := util.Message{
			Message: receivedMessage.Message,
			Author:  username,
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
