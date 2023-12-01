package main

import (
	"errors"
	"fmt"

	"net/http"
	"os"
	"time"

	"turmsapi/router"
)

func main() {

	go router.CheckRooms()

	t := time.Now()
	http.HandleFunc("/createroom", router.CreateRoom)
	http.HandleFunc("/joinroom", router.JoinRoom)

	fmt.Printf("-> Server Running on Port: %v\n-> %v\n", 8080, t.Format(time.UnixDate))
	err := http.ListenAndServe(":"+fmt.Sprint(8080), nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server Closed")
	} else if err != nil {
		fmt.Printf("Error while starting server: %v\n", err)
		os.Exit(1)
	}

}
