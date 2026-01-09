package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/swokamoto/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	

	// CREATE QUERIES OBJECT AFTER REOPENING CONNECTION
	dbQueries := database.New(dbConn)  // ← Move this line here!

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,  // Now using the correct connection
		platform:       platform,
	}

	// Force a new connection
	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Cannot ping database: %v", err)
	}

	// Try multiple ways to check for chirps table
	log.Println("=== Database Debug ===")

	rows, err := dbConn.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' ORDER BY table_name")
	if err != nil {
		log.Printf("Error querying information_schema: %v", err)
	} else {
		log.Println("Tables from information_schema:")
		for rows.Next() {
			var tableName string
			if err := rows.Scan(&tableName); err == nil {
				log.Printf("  - %s", tableName)
			}
		}
		rows.Close()
	}


	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	// mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)


	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpCreate)

	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}