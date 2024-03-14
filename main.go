package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	router := httprouter.New()

	router.POST("/newPassword", createPassword)

	http.ListenAndServe(":8080", router)
}

type data interface {
}

type Password struct {
	Id           int
	Name         string `json:"Name"`
	HashPassword string `json:"HashPassword"`
}

func render(w http.ResponseWriter, password data) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	jContent, _ := json.Marshal(password)
	fmt.Fprintf(w, "%s", jContent)
}

func createPassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var password Password
	err := json.NewDecoder(r.Body).Decode(&password)
	if err != nil {
		log.Println(r.Body)
		http.Error(w, "Could not decode JSON", http.StatusBadRequest)
		return
	}
	conn, _ := sql.Open("sqlite3", "passworby.sdb")
	defer conn.Close()
	_, err = conn.Exec("INSERT INTO Password (Name, HashPassword) VALUES (?, ?)", password.Name, password.HashPassword)
	if err == nil {
		render(w, password)
		log.Println("insert ıs success")
	} else {
		log.Println(err.Error())
		//404 basılabilir
	}
}
