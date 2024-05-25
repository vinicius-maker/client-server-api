package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/vinicius-maker/client-server-api/server/db"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	ApiUrl      = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	EndpointUrl = "/cotacao"
)

type ApiDolarStruct struct {
	Bid string `json:"bid"`
}

func main() {
	conexaoDb := db.NewDb().Conectar()
	defer conexaoDb.Close()

	http.HandleFunc(EndpointUrl, ServerAction(conexaoDb))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}

func ServerAction(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != EndpointUrl {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		cotacaoDolar, err := BuscaCotacaoDolar(ctx)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err = InsertCotacao(ctx, db, cotacaoDolar.Bid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		errJson := json.NewEncoder(w).Encode(cotacaoDolar)
		if errJson != nil {
			log.Printf("erro no encode JSON: %v", errJson)
			return
		}
	}
}

func InsertCotacao(ctx context.Context, db *sql.DB, cotacao string) error {
	_, err := db.ExecContext(ctx, "INSERT INTO cotacoes_dolar (cotacao) VALUES (?)", cotacao)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("erro ao inserir cotacao por dealine do contexto: %w", err)
		}
		return fmt.Errorf("erro ao inserir cotacao no banco de dados: %w", err)
	}

	return nil
}

func BuscaCotacaoDolar(ctx context.Context) (*ApiDolarStruct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", ApiUrl, nil)
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

	var apiDolarStruct struct {
		USDBRL ApiDolarStruct `json:"USDBRL"`
	}

	if err := json.Unmarshal(body, &apiDolarStruct); err != nil {
		return nil, err
	}

	return &apiDolarStruct.USDBRL, nil
}
