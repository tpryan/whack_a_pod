package main

import (
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
	r.HandleFunc("/api/color/", color)
	r.HandleFunc("/api/color-complete", colorComplete)
	r.HandleFunc("/api/color-complete/", colorComplete)
	r.HandleFunc("/color", color)
	r.HandleFunc("/color/", color)
	r.HandleFunc("/color-complete", colorComplete)
	r.HandleFunc("/color-complete/", colorComplete)

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

func colorComplete(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	msg := ""

	h, err := os.Hostname()
	if err != nil {
		msg = fmt.Sprintf("{\"error\":\"could retrieve hostname: %v\"}", err)
		status = http.StatusInternalServerError
	} else {
		msg = fmt.Sprintf("{\"color\":\"%s\", \"name\":\"%s\"}", hexColorString(), h)
	}

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprint(w, msg)
}

func hexColorString() string {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(16777215) // = 0xFFFFFF the highest hex color value allowed.
	return "#" + fmt.Sprintf("%06s", strconv.FormatInt(int64(i), 16))
}
