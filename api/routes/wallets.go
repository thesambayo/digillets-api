package routes

import (
	"errors"
	"net/http"

	"github.com/thesambayo/digillets-api/api/contexts"
	"github.com/thesambayo/digillets-api/api/httpx"
	"github.com/thesambayo/digillets-api/internal/constants"
	"github.com/thesambayo/digillets-api/internal/data/wallets"
	"github.com/thesambayo/digillets-api/internal/publicid"
	"github.com/thesambayo/digillets-api/internal/validators"
)

func (routes *Routes) CreateWallet(resWriter http.ResponseWriter, req *http.Request) {
	user := contexts.ContextGetUser(req)

	var input struct {
		Currency string `json:"currency"`
	}

	err := routes.httpx.ReadJSON(resWriter, req, &input)
	if err != nil {
		routes.httpx.BadRequestResponse(resWriter, req, err)
		return
	}

	validator := validators.New()
	validator.Check(len(input.Currency) != 0, "currency", "currency is required e.g NGN, USD, EUR")
	if !validator.Valid() {
		routes.httpx.FailedValidationResponse(resWriter, req, validator.Errors)
		return
	}

	publicID, _ := publicid.New(constants.PrefixWalletID)

	currency, err := routes.models.Currencies.GetCurrencyByCode(input.Currency)
	if err != nil {
		routes.httpx.ServerErrorResponse(resWriter, req, err)
		return
	}

	wallet := &wallets.Wallet{
		PublicID: publicID,
		Currency: *currency,
		User:     *user,
	}

	wallet, err = routes.models.Wallets.Insert(wallet)

	if err != nil {
		switch {
		case errors.Is(err, constants.ErrDuplicateUserWallet):
			validator.AddError("currency", "a wallet with this currency already exist")
			routes.httpx.FailedValidationResponse(resWriter, req, validator.Errors)
		default:
			routes.httpx.ServerErrorResponse(resWriter, req, err)
		}
		return
	}

	err = routes.httpx.WriteJSON(
		resWriter,
		http.StatusOK,
		httpx.Envelope{"message": "wallet created successfully", "data": wallet},
		nil,
	)

	if err != nil {
		routes.httpx.ServerErrorResponse(resWriter, req, err)
	}
}

func (routes *Routes) GetWallets(resWriter http.ResponseWriter, req *http.Request) {
	user := contexts.ContextGetUser(req)

	wallets, err := routes.models.Wallets.GetByUserId(user.ID)

	if err != nil {
		routes.httpx.ServerErrorResponse(resWriter, req, err)
	}

	err = routes.httpx.WriteJSON(
		resWriter,
		http.StatusOK,
		httpx.Envelope{"message": "wallets fetched successfully", "data": wallets},
		nil,
	)

	if err != nil {
		routes.httpx.ServerErrorResponse(resWriter, req, err)
	}
}

func (routes *Routes) GetSingleWallet(resWriter http.ResponseWriter, req *http.Request) {
	user := contexts.ContextGetUser(req)
	currencyCode := routes.httpx.ReadIDParam(req)

	wallets, err := routes.models.Wallets.GetByCurrencyAndUserId(user.ID, currencyCode)

	if err != nil {
		routes.httpx.ServerErrorResponse(resWriter, req, err)
		return
	}

	err = routes.httpx.WriteJSON(
		resWriter,
		http.StatusOK,
		httpx.Envelope{"message": "wallet fetched successfully", "data": wallets},
		nil,
	)

	if err != nil {
		routes.httpx.ServerErrorResponse(resWriter, req, err)
	}
}
