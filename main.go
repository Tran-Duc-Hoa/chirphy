package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Tran-Duc-Hoa/chirphy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error connecting to the database: %v\n", err)
		return
	}
	defer db.Close()
	dbQueries := database.New(db)

	mux := http.NewServeMux()
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	cfg := &apiConfig{
		db: dbQueries,
		platform: os.Getenv("PLATFORM"),
	}

	mux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {
		var requestBody struct {
			Email string `json:"email"`
		}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&requestBody)
		if err != nil || requestBody.Email == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"error": "Invalid request body"}`)
			return
		}

		user, err := cfg.db.CreateUser(r.Context(), requestBody.Email)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"error": "Failed to create user"}`)
			return
		}

		type responseBody struct {
			ID        uuid.UUID  `json:"id"`
			Email     string `json:"email"`
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
		}

		respBody := responseBody{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}

		data, err := json.Marshal(&respBody)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	})

	mux.HandleFunc("POST /api/chirps", func(w http.ResponseWriter, r *http.Request) {
		var requestBody struct {
			Body string `json:"body"`
			UserId uuid.UUID `json:"user_id"`
		}
		
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&requestBody)
		if err != nil || requestBody.Body == "" || requestBody.UserId == uuid.Nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"error": "Invalid request body"}`)
			return
		}

		chirpyParams := database.CreateChirpyParams{
			Body:   requestBody.Body,
			UserID: requestBody.UserId,
		}
		chirpy, err := cfg.db.CreateChirpy(r.Context(), chirpyParams)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"error": "Failed to create chirp"}`)
			return
		}

		respBody := struct {
			ID        uuid.UUID  `json:"id"`
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
			Body      string `json:"body"`
			UserId    uuid.UUID `json:"user_id"`
		} {
			ID:        chirpy.ID,
			CreatedAt: chirpy.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: chirpy.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			Body:      chirpy.Body,
			UserId:    chirpy.UserID,
		}
		data, err := json.Marshal(&respBody)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	})

	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.Handle("GET /api/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		htmlContent := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())

		fmt.Fprint(w, htmlContent)
	})
	mux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) {
		if cfg.platform != "dev" {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("You are not allowed to reset the hits\n"))
			return
		}
		cfg.fileserverHits.Store(0)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hits reset to 0\n"))
	})

	fmt.Print("Starting server on http://localhost:8080\n")
	httpServer.ListenAndServe()
}

