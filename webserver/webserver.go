package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-redis/redis"
)

var firstDB, secondDB *redis.Client

type data struct {
	Key string
	Val string
}

func main() {

	firstDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer firstDB.Close()
	secondDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})
	defer secondDB.Close()

	http.HandleFunc("/solved", solvedHandler)
	http.HandleFunc("/unsolved", unsolvedHandler)

	log.Fatal(http.ListenAndServe("localhost:80", nil))

}

func solvedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	keys := firstDB.Keys("*")

	var solved []data
	var item data

	for _, key := range keys.Val() {
		item.Key = key
		val, err := firstDB.Get(key).Result()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		item.Val = val

		solved = append(solved, item)
	}

	err := json.NewEncoder(w).Encode(solved)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func unsolvedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	keys := secondDB.Keys("*")

	var solved []data
	var item data

	for _, key := range keys.Val() {
		item.Key = key
		val, err := secondDB.Get(key).Result()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		item.Val = val

		solved = append(solved, item)
	}

	err := json.NewEncoder(w).Encode(solved)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
