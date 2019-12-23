package main

import (
	"log"
	"net/http"
)

func main() {
	r := NewRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	log.Println("start")
	http.Handle("/room", r)
	go r.run()
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
