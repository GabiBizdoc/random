package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	var webRoot string
	if len(os.Args) > 1 {
		webRoot = os.Args[1]
	} else {
		webRoot = "./www"
	}
	fs := noCacheFileServer(http.Dir(webRoot))
	//fs := http.FileServer(http.Dir(webRoot))
	http.Handle("/", fs)

	log.Printf("Serving files from %v on http://localhost:3000", webRoot)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func noCacheFileServer(dir http.Dir) http.Handler {
	fs := http.FileServer(dir)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", time.Now().UTC().Format(http.TimeFormat))
		fs.ServeHTTP(w, r)
	})
}
