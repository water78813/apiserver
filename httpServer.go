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

type Score struct {
	CommentID    string `json:"comment_id`
	ScoreType    string `json:"type`
	CommentScore string `json:"value"`
}

type Data struct {
	Id           string  `json:"id"`
	User         User    `json:"user"`
	Date         string  `json:"date"`
	Text         string  `json:"text"`
	CommentScore []Score `json:"comment_score"`
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
	commentQuery := fmt.Sprintf("SELECT comment_id,comment_date,comment FROM user_comment where comment_id=%s", id)
	commentRows, err := db.Query(commentQuery)
	if err != nil {
		panic(err)
	}
	var comment_id, comment_date, comment string
	for commentRows.Next() {
		err = commentRows.Scan(&comment_id, &comment_date, &comment)
		if err != nil {
			panic(err)
		}
	}
	userData := getData(id)
	user := User{}
	json.Unmarshal(userData, &user)
	d := Response{
		RData: Data{
			Id:           comment_id,
			User:         user,
			Date:         comment_date,
			Text:         comment,
			CommentScore: []Score{},
		},
	}
	scoreQuery := fmt.Sprintf("SELECT score_type,score FROM comment_score where comment_id=%s", id)
	scoreRows, err := db.Query(scoreQuery)
	if err != nil {
		panic(err)
	}
	var score_type, score string
	for scoreRows.Next() {
		err = scoreRows.Scan(&score_type, &score)
		if err != nil {
			panic(err)
		}
		s := Score{
			CommentID:    comment_id,
			ScoreType:    score_type,
			CommentScore: score,
		}
		d.RData.CommentScore = append(d.RData.CommentScore, s)
	}
}
