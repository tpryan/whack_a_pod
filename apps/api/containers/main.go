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

var srv = http.Server{
	ReadTimeout:  5 * time.Second,
	WriteTimeout: 10 * time.Second,
	Addr:         ":8080",
	Handler:      handler(),
}

func main() {
	log.Printf("starting whack a pod color api")
	srv.ListenAndServe()
}

func handler() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/", health)
	r.HandleFunc("/healthz", health)
	r.HandleFunc("/api/healthz", health)
	r.HandleFunc("/api/color", color)
	r.HandleFunc("/api/color-complete", colorComplete)
	r.HandleFunc("/api/color/", color)
	r.HandleFunc("/api/color-complete/", colorComplete)
	return r
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

func color(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, hexColorString())
}

type result struct {
	Color string `json:"color"`
	Name  string `json:"name"`
}

func colorComplete(w http.ResponseWriter, r *http.Request) {
	h, _ := os.Hostname()
	result := result{Name: h, Color: hexColorString()}

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	b, err := json.Marshal(result)
	if err != nil {
		msg := fmt.Sprintf("{\"error\":\"could not unmarshal data %v\"}", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, msg)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(b))
	return
}

func hexColorString() string {
	result := "#"
	for i := 1; i <= 3; i++ {
		rand.Seed(time.Now().UnixNano())
		i := rand.Intn(256)
		result += fmt.Sprintf("%02s", strconv.FormatInt(int64(i), 16))
	}
	return result
}
