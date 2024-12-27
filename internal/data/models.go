package data

import (
	"database/sql"

	"github.com/thesambayo/digillet-api/internal/data/users"
)

type Models struct {
	Users users.UserModel
}

func New(db *sql.DB) *Models {
	return &Models{
		Users: users.UserModel{DB: db},
	}
}

// Create a helper function which returns a Models instance containing the mock models only.
// func NewMockModels() Models {
// 	return Models{
// 		Tickets: MockTicketModel{},
// 		Staff: MockStaffModel{},
// 	}
// }
