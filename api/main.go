package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var allowedOrigins = map[string]bool{
	"http://localhost:5173": true,
	"http://localhost:5174": true,
}

func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./cases.db"
	}
	secret := os.Getenv("TOKEN_SECRET")
	if secret == "" {
		log.Fatal("TOKEN_SECRET is required (set it in .env)")
	}
	portalOrigin := os.Getenv("PORTAL_ORIGIN")
	if portalOrigin == "" {
		portalOrigin = "http://localhost:5173"
	}

	db := openDB(dbPath)
	defer db.Close()

	s := &server{db: db, tokenSecret: secret, portalOrigin: portalOrigin}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/cases", withCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.createCase(w, r)
		case http.MethodGet:
			s.listCases(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	mux.HandleFunc("/api/cases/expire", withCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		s.expireCase(w, r)
	}))
	mux.HandleFunc("/api/cases/delete", withCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		s.deleteCase(w, r)
	}))
	mux.HandleFunc("/api/case", withCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		s.getCase(w, r)
	}))

	log.Printf("api listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
