package currencies

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/thesambayo/digillets-api/internal/constants"
)

type Currency struct {
	ID           int64          `json:"-"`
	Code         string         `json:"code"`
	Name         string         `json:"name"`
	Symbol       string         `json:"symbol"`
	ExchangeRate float64        `json:"-"`
	BaseCurrency sql.NullString `json:"-"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type CurrencyExchangeRate struct {
	Currency       string  `json:"currency"`
	CurrencyName   string  `json:"currency_name"`
	CurrencySymbol string  `json:"currency_symbol"`
	BuyingRate     float64 `json:"buying_rate"`
	SellingRate    float64 `json:"selling_rate"`
}

type CurrencyModel struct {
	DB *sql.DB
}

func (currencyModel *CurrencyModel) GetCurrencyByCode(code string) (*Currency, error) {
	query := fmt.Sprintf(`
		SELECT
			currencies.id,
			currencies.code,
			currencies.name,
			currencies.symbol
		FROM
			currencies
		WHERE
			currencies.code = '%v';
	`, code)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var currency Currency
	err := currencyModel.DB.QueryRowContext(ctx, query).Scan(
		&currency.ID,
		&currency.Code,
		&currency.Name,
		&currency.Symbol,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, constants.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &currency, err
}

func (currencyModel *CurrencyModel) GetExchangeRateBetweenTwoCurrencies(currFrom, currTo string) (float64, error) {
	query := fmt.Sprintf(`
    SELECT
      ROUND((currencyTo.exchange_rate / currencyFrom.exchange_rate), 4) AS conversion_rate
    FROM
      currencies currencyFrom
    JOIN
      currencies currencyTo ON currencyTo.code = '%v'
    WHERE
      currencyFrom.code = '%v';
  `, currTo, currFrom)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var conversionRate float64
	err := currencyModel.DB.QueryRowContext(ctx, query).Scan(&conversionRate)
	return conversionRate, err
}

func (currencyModel *CurrencyModel) GetExchangeRatesForACurrency(currency string) ([]*CurrencyExchangeRate, error) {
	query := `
		WITH base_currency AS (
	    SELECT exchange_rate AS rate
	    FROM currencies
	    WHERE code = $1
		)
		SELECT
	    currency.code AS currency,
	    currency.name AS currency_name,
			currency.symbol AS currency_symbol,
	    ROUND(base_currency.rate / currency.exchange_rate, 6) AS buying_rate,
	    ROUND(currency.exchange_rate / base_currency.rate, 6) AS selling_rate
		FROM
	    currencies currency, base_currency
		WHERE
	    currency.code != $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{
		currency,
	}
	rows, err := currencyModel.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var exchangeRates []*CurrencyExchangeRate
	for rows.Next() {
		var exchangeRate CurrencyExchangeRate

		err := rows.Scan(
			&exchangeRate.Currency,
			&exchangeRate.CurrencyName,
			&exchangeRate.CurrencySymbol,
			&exchangeRate.BuyingRate,
			&exchangeRate.SellingRate,
		)

		if err != nil {
			return nil, err
		}
		exchangeRates = append(exchangeRates, &exchangeRate)
	}
	return exchangeRates, nil
}
