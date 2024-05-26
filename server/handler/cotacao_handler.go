package handler

import (
	"encoding/json"
	"github.com/vinicius-maker/client-server-api/server/service"
	"log"
	"net/http"
)

const EndpointCotacaoUrl = "/cotacao"

type CotacaoHandler struct {
	BuscarCotacaoDolarService service.BuscarCotacaoDolar
}

func NewCotacaoHandler(buscarCotacaoService service.BuscarCotacaoDolar) *CotacaoHandler {
	return &CotacaoHandler{
		BuscarCotacaoDolarService: buscarCotacaoService,
	}
}

func (c *CotacaoHandler) HandlerCotacaoDolar(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != EndpointCotacaoUrl {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cotacaoDolar, err := c.BuscarCotacaoDolarService.BuscaCotacaoDolar()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if errJson := json.NewEncoder(w).Encode(cotacaoDolar); errJson != nil {
		log.Printf("erro no encode JSON: %v", errJson)
		return
	}
}
