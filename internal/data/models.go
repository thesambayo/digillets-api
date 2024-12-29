package data

import (
	"database/sql"

	"github.com/thesambayo/digillets-api/internal/data/currencies"
	"github.com/thesambayo/digillets-api/internal/data/users"
	"github.com/thesambayo/digillets-api/internal/data/wallets"
)

type Models struct {
	Users      users.UserModel
	Currencies currencies.CurrencyModel
	Wallets    wallets.WalletModel
}

func New(db *sql.DB) *Models {
	return &Models{
		Users:      users.UserModel{DB: db},
		Currencies: currencies.CurrencyModel{DB: db},
		Wallets:    wallets.WalletModel{DB: db},
	}
}

// Create a helper function which returns a Models instance containing the mock models only.
// func NewMockModels() Models {
// 	return Models{
// 		Tickets: MockTicketModel{},
// 		Staff: MockStaffModel{},
// 	}
// }
