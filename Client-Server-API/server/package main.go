package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type DollarQuote struct {
	Usdbrl struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/", CotacaoUsdHandleer)
	http.ListenAndServe(":8080", nil)
}

func CotacaoUsdHandleer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	dollarQuote, err := CotacaoUsd(ctx)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	select {
	case <-ctx.Done():
		log.Println("Tempo limite excedido")
		http.Error(w, "Tempo limite excedido", http.StatusGatewayTimeout)
		return
	default:
		log.Printf("Cotação do dólar hoje: %s\n", dollarQuote.Usdbrl.Bid)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(dollarQuote.Usdbrl)
	}
}

func CotacaoUsd(ctx context.Context) (*DollarQuote, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var dollarQuote DollarQuote
	err = json.Unmarshal(body, &dollarQuote)
	if err != nil {
		return nil, err
	}
	return &dollarQuote, nil
}
