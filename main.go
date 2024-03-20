package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	DB             *DB
	JWTSecret      string
}

func main() {
	godotenv.Load("key.env")
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	//if --debug is set, delete database.json
	if *dbg {
		fmt.Println("trying to delete old database...")
		err := delDB()
		if err != nil {
			fmt.Println(err)
		}
	}

	db, err := NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
		JWTSecret:      os.Getenv("JWT_SECRET"),
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
	mux.HandleFunc("POST /api/login", apiCfg.loginHandler)
	//mux.HandleFunc("PUT /api/users", apiCfg.updateUser)

	corsMux := middlewareCors(mux)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on %s\n", http.Dir("."), srv.Addr)

	log.Fatal(srv.ListenAndServe())
}
