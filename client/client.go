package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	ServerUrl = "http://localhost:8080/cotacao"
	FileName  = "cotacao.txt"
)

type ApiDolarStruct struct {
	Cotacao string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	cotacaoDolar, err := BuscaCotacaoDolar(ctx)
	if err != nil {
		log.Fatalf("erro ao buscar a cotação do dólar: %v", err)
	}

	err = CriaArquivo(cotacaoDolar.Cotacao)
	if err != nil {
		log.Fatalf("erro ao criar o arquivo: %v", err)
	}
}

func BuscaCotacaoDolar(ctx context.Context) (*ApiDolarStruct, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", ServerUrl, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var apiDolarStruct ApiDolarStruct
	if errJson := json.Unmarshal(body, &apiDolarStruct); errJson != nil {
		return nil, errJson
	}

	return &apiDolarStruct, nil
}

func CriaArquivo(dolar string) error {
	arquivo, err := os.Create(FileName)
	if err != nil {
		return err
	}
	defer arquivo.Close()

	_, err = arquivo.WriteString("Dólar: " + dolar)
	return err
}
