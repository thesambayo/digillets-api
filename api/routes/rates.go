package routes

import (
	"net/http"

	"github.com/thesambayo/digillets-api/api/httpx"
	"github.com/thesambayo/digillets-api/internal/validators"
)

func (routes *Routes) GetTwoCurrenciesExchangeRate(resWriter http.ResponseWriter, req *http.Request) {
	var input struct {
		currencyFrom string
		currencyTo   string
	}

	queryString := req.URL.Query()
	input.currencyFrom = routes.httpx.ReadString(queryString, "currencyFrom", "")
	input.currencyTo = routes.httpx.ReadString(queryString, "currencyTo", "")

	validator := validators.New()
	validator.Check(len(input.currencyFrom) != 0, "currencyFrom", "currencyFrom must have a value e.g NGN, USD, EUR")
	validator.Check(len(input.currencyTo) != 0, "currencyTo", "currencyTo must have a value e.g NGN, USD, EUR")
	if !validator.Valid() {
		routes.httpx.FailedValidationResponse(resWriter, req, validator.Errors)
		return
	}

	rate, err := routes.models.Currencies.GetExchangeRateBetweenTwoCurrencies(input.currencyFrom, input.currencyTo)
	data := struct {
		ExchangeRate float64 `json:"exchange_rate"`
	}{
		ExchangeRate: rate,
	}

	if err != nil {
		routes.httpx.ServerErrorResponse(resWriter, req, err)
		return
	}

	err = routes.httpx.WriteJSON(resWriter, http.StatusOK, httpx.Envelope{"data": data}, nil)
	if err != nil {
		routes.httpx.ServerErrorResponse(resWriter, req, err)
		return
	}
}

func (routes *Routes) GetCurrencyExchangeRates(resWriter http.ResponseWriter, req *http.Request) {
	var input struct {
		currency string
	}

	queryString := req.URL.Query()
	input.currency = routes.httpx.ReadString(queryString, "currency", "")

	validator := validators.New()
	validator.Check(len(input.currency) != 0, "currency", "currency is required e.g NGN, USD, EUR")
	if !validator.Valid() {
		routes.httpx.FailedValidationResponse(resWriter, req, validator.Errors)
		return
	}

	exchangeRates, err := routes.models.Currencies.GetExchangeRatesForACurrency(input.currency)

	if err != nil {
		routes.httpx.ServerErrorResponse(resWriter, req, err)
		return
	}

	err = routes.httpx.WriteJSON(resWriter, http.StatusOK, httpx.Envelope{"data": exchangeRates}, nil)
	if err != nil {
		routes.httpx.ServerErrorResponse(resWriter, req, err)
		return
	}
}
