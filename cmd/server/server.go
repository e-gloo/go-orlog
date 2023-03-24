package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/e-gloo/orlog/orlog/server"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":8080", "http service address")
var upgrader = websocket.Upgrader{} // use default options


func main() {
	flag.Parse()
	hub := server.NewHub()
	go hub.Run()
	http.HandleFunc("/connect", server.MessageHandler)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.ServeWs(hub, w, r)
	})
	fmt.Println("Listening on ", addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
