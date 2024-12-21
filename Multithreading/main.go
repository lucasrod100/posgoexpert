package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type ResponseCEP struct {
	API        string `json:"api"`
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	Uf         string `json:"uf"`
}

type ResponseCEPError struct {
	API   string
	Error string
}

func requestAPICep(ctx context.Context, url string) (http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return http.Response{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return http.Response{}, err
	}

	if res.StatusCode != http.StatusOK {
		return http.Response{}, err
	}

	return *res, nil
}

func requestAPICepViaCEP(ctx context.Context, cep string, ch chan<- ResponseCEP, chError chan<- ResponseCEPError) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	apiName := "ViaCEP"

	res, err := requestAPICep(ctx, url)
	if err != nil {
		chError <- ResponseCEPError{API: apiName, Error: err.Error()}
		return
	}
	defer res.Body.Close()

	var response ResponseCEP
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		chError <- ResponseCEPError{API: apiName, Error: err.Error()}
		return
	}

	response.API = apiName
	ch <- response
}

func requestAPICepBrasilAPI(ctx context.Context, cep string, ch chan<- ResponseCEP, chError chan<- ResponseCEPError) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	apiName := "BrasilAPI"

	res, err := requestAPICep(ctx, url)
	if err != nil {
		chError <- ResponseCEPError{API: apiName, Error: err.Error()}
		return
	}
	defer res.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		chError <- ResponseCEPError{API: apiName, Error: err.Error()}
		return
	}

	response := ResponseCEP{
		Cep:        data["cep"].(string),
		Logradouro: data["street"].(string),
		Bairro:     data["neighborhood"].(string),
		Localidade: data["city"].(string),
		Uf:         data["state"].(string),
		API:        apiName,
	}

	ch <- response
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <CEP>")
		os.Exit(1)
	}

	cep := os.Args[1]

	if len(cep) != 8 {
		fmt.Println("Erro: O CEP deve conter exatamente 8 números.")
		os.Exit(1)
	}

	ch := make(chan ResponseCEP, 2)
	chError := make(chan ResponseCEPError, 2)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	go requestAPICepViaCEP(ctx, cep, ch, chError)
	go requestAPICepBrasilAPI(ctx, cep, ch, chError)

	select {
	case response := <-ch:
		fmt.Printf("Resposta mais rápida da API %s:\n", response.API)
		fmt.Printf("CEP: %s\nLogradouro: %s\nBairro: %s\nLocalidade: %s\nUF: %s\n", response.Cep, response.Logradouro, response.Bairro, response.Localidade, response.Uf)
	case responseError := <-chError:
		fmt.Printf("Erro na API %s: %s\n", responseError.API, responseError.Error)
	case <-ctx.Done():
		fmt.Println("Requisição excedeu o tempo limite de 1 segundo")
	}
}
