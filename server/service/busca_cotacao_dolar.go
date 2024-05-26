package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const ApiUrl = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

type ApiDolarStruct struct {
	Bid string `json:"bid"`
}

type BuscarCotacaoDolar struct {
	Db *sql.DB
}

func NewBuscarCotacaoDolar(db *sql.DB) BuscarCotacaoDolar {
	return BuscarCotacaoDolar{
		Db: db,
	}
}

func (b BuscarCotacaoDolar) BuscaCotacaoDolar() (*ApiDolarStruct, error) {
	cotacaoDolar, err := b.BuscaCotacaoDolarApi()
	if err != nil {
		return cotacaoDolar, err
	}

	err = b.InsertCotacaoDolar(cotacaoDolar.Bid)
	if err != nil {
		return cotacaoDolar, err
	}

	return cotacaoDolar, nil
}

func (b *BuscarCotacaoDolar) BuscaCotacaoDolarApi() (*ApiDolarStruct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", ApiUrl, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("erro ao buscar cotacao por deadline do contexto: %w", err)
		}
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

func (b *BuscarCotacaoDolar) InsertCotacaoDolar(cotacao string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := b.Db.ExecContext(ctx, "INSERT INTO cotacoes_dolar (valor_cotacao) VALUES (?)", cotacao)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("erro ao inserir cotacao por deadline do contexto: %w", err)
		}
		return fmt.Errorf("erro ao inserir cotacao no banco de dados: %w", err)
	}

	return nil
}
