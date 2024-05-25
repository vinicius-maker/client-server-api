package main

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/vinicius-maker/client-server-api/server/db"
	"github.com/vinicius-maker/client-server-api/server/handler"
	"github.com/vinicius-maker/client-server-api/server/service"
	"log"
	"net/http"
)

const EndpointUrl = "/cotacao"

func main() {
	conexaoDb := db.NewDb().Conectar()
	if conexaoDb == nil {
		log.Fatalf("erro ao conectar ao banco de dados")
		return
	}
	defer conexaoDb.Close()

	buscarCotacaoDolarService := service.NewBuscarCotacaoDolar(conexaoDb)
	if buscarCotacaoDolarService.Db == nil {
		log.Fatalf("erro ao criar o serviço de busca de cotação do dólar")
		return
	}

	cotacaoHandler := handler.NewCotacaoHandler(buscarCotacaoDolarService)
	if cotacaoHandler.BuscarCotacaoDolarService.Db == nil {
		log.Fatalf("erro ao criar o handler de cotação do dólar")
		return
	}

	http.HandleFunc(EndpointUrl, cotacaoHandler.HandlerCotacaoDolar)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("erro ao iniciar o servidor: %v", err)
	}
}
