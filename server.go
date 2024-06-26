package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	apiURL        = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	serverAddress = ":8080"
	databaseFile  = "./cotacoes.db"
	apiTimeout    = 200 * time.Millisecond
	dbTimeout     = 10 * time.Millisecond
)

type Cotacao struct {
	Bid string `json:"bid"`
}

type APIResponse struct {
	USDBRL Cotacao `json:"USDBRL"`
}

func fetchCotacao(ctx context.Context) (*Cotacao, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code")
	}

	var apiResponse APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}

	return &apiResponse.USDBRL, nil
}

func saveCotacao(ctx context.Context, db *sql.DB, cotacao *Cotacao) error {
	query := "INSERT INTO cotacoes (bid, timestamp) VALUES (?, ?)"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, cotacao.Bid, time.Now())
	return err
}

func cotacaoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), apiTimeout)
		defer cancel()

		cotacao, err := fetchCotacao(ctx)
		if err != nil {
			http.Error(w, fmt.Sprintf("error fetching cotacao: %v", err), http.StatusInternalServerError)
			log.Println("Error fetching cotacao:", err)
			return
		}

		ctx, cancel = context.WithTimeout(r.Context(), dbTimeout)
		defer cancel()

		if err := saveCotacao(ctx, db, cotacao); err != nil {
			http.Error(w, fmt.Sprintf("error saving cotacao: %v", err), http.StatusInternalServerError)
			log.Println("Error saving cotacao:", err)
			return
		}

		response := map[string]string{"bid": cotacao.Bid}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func setupDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		return nil, err
	}

	query := `
	CREATE TABLE IF NOT EXISTS cotacoes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bid TEXT,
		timestamp DATETIME
	)`
	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	db, err := setupDatabase()
	if err != nil {
		log.Fatalf("error setting up database: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/cotacao", cotacaoHandler(db))
	log.Println("Server started at", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, nil))
}
