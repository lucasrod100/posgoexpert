package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type ResponseDollarQuote struct {
	Bid string
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatalf("Erro ao criar requisição: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)

	select {
	case <-ctx.Done():
		log.Fatalf("Tempo limite excedido para obter resposta da requisição: %v", ctx.Err())
	default:
		if err != nil {
			log.Fatalf("Erro ao fazer requisição: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Erro do servidor: %s", resp.Status)
		}

		var responseDollarQuote ResponseDollarQuote
		err = json.NewDecoder(resp.Body).Decode(&responseDollarQuote)
		if err != nil {
			log.Fatalf("Erro ao obter resposta: %v", err)
		}

		fmt.Println("Cotação atual do dolar: ", responseDollarQuote.Bid)

		saveQuoteFile(responseDollarQuote)
	}
}

func saveQuoteFile(responseDollarQuote ResponseDollarQuote) {
	err := os.WriteFile("cotacao.txt", []byte(fmt.Sprintf("Dólar: %s", responseDollarQuote.Bid)), 0644)
	if err != nil {
		log.Fatalf("Erro ao salvar cotação no arquivo: %v", err)
	}

	fmt.Println("Cotação salva no arquivo cotacao.txt")
}
