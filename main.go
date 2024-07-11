package main

import (
	"log"
	"net/http"
)

func main() {

	const filepathRoot string = "."
	const port string = "8080"

	mux := http.NewServeMux()
	apiCfg := new(apiConfig)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("/api/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	fileSrv := http.FileServer(http.Dir(filepathRoot))
	fileSrv = http.StripPrefix("/app", fileSrv)
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(fileSrv))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
