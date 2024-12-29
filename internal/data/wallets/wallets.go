package wallets

import (
	"context"
	"database/sql"
	"time"

	"github.com/thesambayo/digillets-api/internal/constants"
	"github.com/thesambayo/digillets-api/internal/data/currencies"
	"github.com/thesambayo/digillets-api/internal/data/users"
)

// Wallet represents the wallets table in the database.
type Wallet struct {
	ID        string              `json:"-"`
	PublicID  string              `json:"public_id"`
	User      users.User          `json:"user"`
	Currency  currencies.Currency `json:"currency"`
	Balance   float64             `json:"balance"`
	IsFrozen  bool                `json:"is_frozen"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

type WalletModel struct {
	DB *sql.DB
}

func (walletModel WalletModel) Insert(wallet *Wallet) (*Wallet, error) {

	query := `
    INSERT INTO wallets
      (public_id, user_id, currency_id, balance, is_frozen)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id, created_at`

	args := []interface{}{
		wallet.PublicID,
		wallet.User.ID,
		wallet.Currency.ID,
		wallet.Balance,
		wallet.IsFrozen,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// If the table already contains a record with this email address, then when we try
	// to perform the insert there will be a violation of the UNIQUE "user_email_key"
	// constraint that we set up in the previous chapter. We check for this error
	// specifically, and return custom ErrDuplicateEmail error instead.
	err := walletModel.DB.QueryRowContext(ctx, query, args...).Scan(&wallet.ID, &wallet.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "unique_user_currency"`:
			return nil, constants.ErrDuplicateUserWallet
		default:
			return nil, err
		}
	}

	return wallet, nil
}

func (walletModel WalletModel) GetByUserId(userId int64) ([]*Wallet, error) {
	query := `
		SELECT
			wallets.public_id,
		  wallets.balance,
		  wallets.is_frozen,
		  wallets.created_at,
		  wallets.updated_at,
		  currencies.code AS currency_code,
		  currencies.name AS currency_name,
		  currencies.symbol AS currency_symbol
		FROM
	 		wallets
		JOIN
			currencies ON currencies.id = wallets.currency_id
		WHERE
	  	wallets.user_id = $1;
  `

	args := []interface{}{
		userId,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// If the table already contains a record with this email address, then when we try
	// to perform the insert there will be a violation of the UNIQUE "user_email_key"
	// constraint that we set up in the previous chapter. We check for this error
	// specifically, and return custom ErrDuplicateEmail error instead.
	rows, err := walletModel.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []*Wallet
	for rows.Next() {
		var wallet Wallet
		err := rows.Scan(
			&wallet.PublicID,
			&wallet.Balance,
			&wallet.IsFrozen,
			&wallet.CreatedAt,
			&wallet.UpdatedAt,
			&wallet.Currency.Code,
			&wallet.Currency.Name,
			&wallet.Currency.Symbol,
		)

		if err != nil {
			return nil, err
		}

		wallets = append(wallets, &wallet)
	}

	return wallets, nil
}

func (walletModel WalletModel) GetByCurrencyAndUserId(userId int64, currencyCode string) (*Wallet, error) {
	query := `
		SELECT
			wallets.public_id,
		  wallets.balance,
		  wallets.is_frozen,
		  wallets.created_at,
		  wallets.updated_at,
		  currencies.code AS currency_code,
		  currencies.name AS currency_name,
		  currencies.symbol AS currency_symbol
		FROM
	 		wallets
		JOIN
			currencies ON currencies.id = wallets.currency_id
		WHERE
	  	wallets.user_id = $1
		AND
			currencies.code = $2;
  `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{
		userId,
		currencyCode,
	}

	var wallet Wallet
	err := walletModel.DB.QueryRowContext(ctx, query, args...).Scan(
		&wallet.PublicID,
		&wallet.Balance,
		&wallet.IsFrozen,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
		&wallet.Currency.Code,
		&wallet.Currency.Name,
		&wallet.Currency.Symbol,
	)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}
