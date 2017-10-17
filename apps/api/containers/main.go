package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	log.Printf("starting whack a pod api api")
	http.ListenAndServe(":8080", handler())
}

func handler() http.Handler {

	r := http.NewServeMux()
	r.HandleFunc("/", health)
	r.HandleFunc("/healthz", health)
	r.HandleFunc("/api/healthz", health)
	r.HandleFunc("/api/color", color)
	r.HandleFunc("/api/color-complete", colorComplete)

	return r
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

func color(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, hexColorString())
}

type result struct {
	Color string `json:"color"`
	Name  string `json:"name"`
}

func colorComplete(w http.ResponseWriter, r *http.Request) {
	h, _ := os.Hostname()
	result := result{
		Name:  h,
		Color: hexColorString(),
	}

	b, err := json.Marshal(result)
	if err != nil {
		msg := fmt.Sprintf("{\"error\":\"could not unmarshap data %v\"}", err)
		sendJSON(w, msg, http.StatusInternalServerError)
	}

	sendJSON(w, string(b), http.StatusOK)
}

func hexColorString() string {
	result := "#"
	for i := 1; i <= 3; i++ {
		rand.Seed(time.Now().UnixNano())
		i := rand.Intn(256)
		result += strconv.FormatInt(int64(i), 16)
	}
	return result
}

func sendJSON(w http.ResponseWriter, content string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprint(w, content)
}
