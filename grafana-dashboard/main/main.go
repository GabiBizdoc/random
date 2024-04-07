package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	log.Println("starting the server...")
	appName := strings.ToUpper(os.Getenv("APP_NAME"))
	switch appName {
	case "HTTP_SERVER":
		StartServer(7000)
	case "HTTP_SERVER_BTREE":
		go doWork()
		StartServer(7000)
	default:
		err := fmt.Errorf("invalid value for APP_NAME (HTTP_SERVER, HTTP_SERVER_BTREE)")
		log.Println(err)
		panic(err)
	}
}

func StartServer(port int) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"hello": "world", "now": "%s"}`, time.Now().Format(time.RFC3339))
	})
	host := fmt.Sprintf("0.0.0.0:%d", port)
	log.Println("http://" + host)
	log.Println(http.ListenAndServe(host, mux))
	log.Println("closed")
}
