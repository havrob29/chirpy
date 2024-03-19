package main

import (
	"flag"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
	DB             *DB
}

func main() {
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *dbg {

	}
	db, err := NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}
	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/*", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/admin/metrics", apiCfg.adminMetrics)
	mux.HandleFunc("/api/healthz", handlerReadiness)
	mux.HandleFunc("/api/reset", apiCfg.resetNumberRequests)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerSingleRetrieve)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUserCreate)

	corsMux := middlewareCors(mux)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on %s\n", http.Dir("."), srv.Addr)

	log.Fatal(srv.ListenAndServe())
}
