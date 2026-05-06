package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("harbor up!\n"))
	})

	log.Fatal(http.ListenAndServe(":10007", nil))
}
