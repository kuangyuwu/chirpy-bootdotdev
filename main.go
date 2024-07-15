package main

import (
	"log"
	"net/http"

	"github.com/kuangyuwu/chirpy-bootdev/internal/database"
)

func main() {

	const filepathRoot string = "."
	const port string = "8080"

	mux := http.NewServeMux()
	dbPath := "./database.json"
	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	apiCfg := &apiConfig{
		fileserverHits: 0,
		db:             db,
	}

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("/api/reset", apiCfg.handlerReset)
	// mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerPostChirp)

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
