package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("POST")
	r.HandleFunc("/v1/comments/", getCommentsIDHandler).Methods("GET")
	r.HandleFunc("/v2/comments/{id}", getDataHandler).Methods("GET")
	server := &http.Server{
		Handler: r,
		Addr:    ":8080",
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func homeHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		fmt.Println("it is POST")
	}
	w.WriteHeader(200)
	fmt.Println("here")
	s := "welcome"
	w.Write([]byte(s))
}

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Data struct {
	Id   string `json:"id"`
	User User   `json:"user"`
	Date string `json:"date"`
	Text string `json:"text"`
}

type Response struct {
	RData Data `json:"data"`
}

func getCommentsIDHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if req.Method == "GET" {
		var d = Data{
			Id:   "1",
			User: User{},
			Date: "2015-01-01",
			Text: "左手には少しさがって博物の教室がある。\r",
		}
		json.NewEncoder(w).Encode(d)
	}
}
func getDataHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	p := mux.Vars(req)
	id := p["id"]
	d := accessDB(id)
	json.NewEncoder(w).Encode(d)
}

func getData(id string) []byte {
	url := fmt.Sprintf("https://interview-external.moneywelfare.com/users/%s", id)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return d
}

func accessDB(id string) Response {
	db, err := sql.Open("postgres", "user=postgres password=password dbname=commentsdb sslmode=disable")
	defer db.Close()
	if err != nil {
		panic(err)
	}
	q := fmt.Sprintf("SELECT * FROM user_comment where comment_id=%s", id)
	rows, err := db.Query(q)
	if err != nil {
		panic(err)
	}
	var comment_id, user_id, comment_date, comment string
	for rows.Next() {
		err = rows.Scan(&comment_id, &user_id, &comment_date, &comment)
		if err != nil {
			panic(err)
		}
		fmt.Println(comment_id, user_id, comment_date, comment)
	}
	userData := getData(id)
	user := User{}
	json.Unmarshal(userData, &user)
	d := Response{
		RData: Data{
			Id:   comment_id,
			User: user,
			Date: comment_date,
			Text: comment,
		},
	}
	return d
}
