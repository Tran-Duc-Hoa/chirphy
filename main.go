package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	cfg := &apiConfig{}

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
		cfg.fileserverHits.Store(0)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hits reset to 0\n"))
	})
	mux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {
		var requestBody struct {
			Body string `json:"body"`
		}

		 decoder:= json.NewDecoder(r.Body)
		err := decoder.Decode(&requestBody)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"error": "Something went wrong"}`)
			return
		}

		if len(requestBody.Body) > 140 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"error": "Chirp is too long"}`)
			return
		}

		words := strings.Split(requestBody.Body, " ")
		for i, word := range words {
			loweredWord := strings.ToLower(word)
			if loweredWord == "kerfuffle" || loweredWord == "sharbert" || loweredWord == "fornax" {
				words[i] = "****"
			}
		}

		type responseBody struct {
			Valid   bool   `json:"valid"`
			CleanedBody   string `json:"cleaned_body"`
		}
		respBody := responseBody{
			Valid: true,
			CleanedBody: strings.Join(words, " "),
		}
		data, err := json.Marshal(&respBody)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})

	fmt.Print("Starting server on http://localhost:8080\n")
	httpServer.ListenAndServe()
}

