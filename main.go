package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Todo struct {
	Id          string `json:id`
	Title       string `json:"title"`
	Desc        string `json:"desc"`
	IsCompleted bool   `json:"isCompleted"`
}

var todos []Todo

func main() {
	fmt.Println("start")
	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Authorization"})

	todos = []Todo{
		Todo{Id: "1", Title: "learn Next.js", Desc: "Description", IsCompleted: false},
		Todo{Id: "2", Title: "learn tailwind", Desc: "Description", IsCompleted: false},
	}
	r := mux.NewRouter()
	r.Handle("/public", http.HandlerFunc(getPublic)).Methods("GET")
	r.Handle("/todos", http.HandlerFunc(getTodos)).Methods("GET")
	http.ListenAndServe(":8080", handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(r))
}

func getPublic(w http.ResponseWriter, r *http.Request) {
	paylaod, _ := json.Marshal("public api")
	fmt.Fprint(w, string(paylaod))
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	paylaod, _ := json.Marshal(todos)
	fmt.Fprint(w, string(paylaod))
}
