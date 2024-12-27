package routes

import (
	"net/http"

	"github.com/thesambayo/digillet-api/api/httpx"
)

func (routes *Routes) HealthcheckHandler(resWriter http.ResponseWriter, req *http.Request) {
	data := httpx.Envelope{
		"status": "available",
	}

	err := routes.httpx.WriteJSON(resWriter, http.StatusOK, data, nil)
	if err != nil {
		routes.httpx.ServerErrorResponse(resWriter, req, err)
		return
	}
}
