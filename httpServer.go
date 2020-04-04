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

type Data struct {
	Id   int    `json:"id"`
	User int    `json:"user_id"`
	Date string `json:"date"`
	Text string `json:"text"`
}

func getCommentsIDHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if req.Method == "GET" {
		var d = Data{
			Id:   1,
			User: 2,
			Date: "2015-01-01",
			Text: "左手には少しさがって博物の教室がある。\r",
		}
		json.NewEncoder(w).Encode(d)
	}
}
