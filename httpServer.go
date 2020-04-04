package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/v1/comments/:id", getCommentsIDHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func homeHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		fmt.Println("it is POST")
	}
	w.WriteHeader(200)
	fmt.Println("here")
}

type data struct {
	id     int    `json:"id"`
	userID int    `json:"user_id"`
	date   string `json:"date"`
	text   string `json:"text"`
}

var d = &data{
	id:     1,
	userID: 2,
	date:   "2015-01-01",
	text:   "左手には少しさがって博物の教室がある。\r",
}

func getCommentsIDHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		j, err := json.Marshal(d)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(j)
		fmt.Println("get")
	}
}
