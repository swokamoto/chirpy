package main

import (
	"net/http"
	"github.com/joho/godotenv"
	"os"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	godotenv.Load()
	if os.Getenv("PLATFORM") != "dev" {
		http.Error(w, "403 Forbidden", http.StatusForbidden)
		return
	}
	
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
