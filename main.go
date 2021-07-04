package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go"
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

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// config := &firebase.Config{ProjectID: "	golang-todo-app"}
		app, err := firebase.NewApp(context.Background(), nil)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
		auth, err := app.Auth(context.Background())
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
		authHeader := r.Header.Get("Authorization")
		idToken := strings.Replace(authHeader, "Bearer", "", 1)

		token, err := auth.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			// JWT が無効なら Handler に進まず別処理
			fmt.Printf("error verifying ID token: %v\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("error verifying ID token\n"))
			return
		}
		log.Printf("Verified ID token: %v\n", token)
		next.ServeHTTP(w, r)
	}
}

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
	r.Handle("/todos", authMiddleware(http.HandlerFunc(getTodos))).Methods("GET")
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
