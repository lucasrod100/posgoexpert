package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DollarQuote struct {
	Usdbrl struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

var db *sql.DB

func main() {
	setupDB()
	defer db.Close()
	http.HandleFunc("/", CotacaoUsdHandleer)
	http.ListenAndServe(":8080", nil)
}

func CotacaoUsdHandleer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	resp, err := CotacaoUsd(ctx)

	select {
	case <-ctx.Done():
		log.Println("Tempo limite excedido")
		http.Error(w, "Tempo limite excedido", http.StatusGatewayTimeout)
		return
	default:
		defer resp.Body.Close()
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var dollarQuote DollarQuote
		err = json.Unmarshal(body, &dollarQuote)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = saveDollarQuote(dollarQuote)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Cotação do dólar hoje: %s\n", dollarQuote.Usdbrl.Bid)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(dollarQuote.Usdbrl)
	}
}

func CotacaoUsd(ctx context.Context) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req)
}

func setupDB() {
	var err error
	db, err = sql.Open("sqlite3", "./cotacoes.db")
	if err != nil {
		log.Fatalf("Erro ao abrir o banco de dados: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cotacoes (id INTEGER PRIMARY KEY AUTOINCREMENT, bid TEXT, date DATETIME DEFAULT CURRENT_TIMESTAMP)`)
	if err != nil {
		log.Fatalf("Erro ao criar tabela no banco de dados: %v", err)
	}
}

func saveDollarQuote(dollarQuote DollarQuote) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// Inserção no banco de dados
	query := `INSERT INTO cotacoes (bid) VALUES (?)`
	_, err := db.ExecContext(ctx, query, dollarQuote.Usdbrl.Bid)

	select {
	case <-ctx.Done():
		return fmt.Errorf("operação para inserir a cotação excedeu o tempo limite: %w", ctx.Err())
	default:
		if err != nil {
			return fmt.Errorf("erro ao inserir a cotação: %w", err)
		}
		return nil
	}
}
