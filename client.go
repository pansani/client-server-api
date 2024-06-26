package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	serverURL     = "http://localhost:8080/cotacao"
	clientTimeout = 300 * time.Millisecond
	outputFile    = "cotacao.txt"
)

func fetchCotacao(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", serverURL, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response code")
	}

	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response["bid"], nil
}

func saveCotacaoToFile(cotacao string) error {
	content := fmt.Sprintf("Dólar: %s", cotacao)
	return ioutil.WriteFile(outputFile, []byte(content), 0644)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), clientTimeout)
	defer cancel()

	cotacao, err := fetchCotacao(ctx)
	if err != nil {
		log.Fatalf("error fetching cotacao: %v", err)
	}

	if err := saveCotacaoToFile(cotacao); err != nil {
		log.Fatalf("error saving cotacao to file: %v", err)
	}

	log.Println("Cotacao saved to file:", outputFile)
}
