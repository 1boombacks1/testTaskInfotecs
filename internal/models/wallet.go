package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Wallet struct {
	ID      uuid.UUID `db:"id"`
	Balance float32   `db:"balance"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
