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
	polka_key      string
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
		polka_key:      os.Getenv("POLKA_KEY"),
	}

	mux := http.NewServeMux()
	mux.Handle("/app/*", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/admin/metrics", apiCfg.adminMetrics)
	mux.HandleFunc("/api/healthz", handlerReadiness)
	mux.HandleFunc("/api/reset", apiCfg.resetNumberRequests)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.getApiChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getApiChirpsByID)
	mux.HandleFunc("POST /api/users", apiCfg.postApiUsers)
	mux.HandleFunc("POST /api/login", apiCfg.postApiLogin)
	mux.HandleFunc("PUT /api/users", apiCfg.putApiUser)
	mux.HandleFunc("POST /api/refresh", apiCfg.postApiRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.postApiRevoke)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.deleteChirp)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.postApiPolkaWebhook)

	corsMux := middlewareCors(mux)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on %s\n", http.Dir("."), srv.Addr)

	log.Fatal(srv.ListenAndServe())
}
