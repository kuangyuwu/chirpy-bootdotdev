package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	const port string = "8080"

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
