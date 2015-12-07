package main

import (
	"io"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "golang https server")
}

func main() {
	http.HandleFunc("/", handler)
	if err := http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
