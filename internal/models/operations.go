package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Operations struct {
	ID        int       `db:"id" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"time"`

	FromWalletID *uuid.UUID `db:"from_wallet_id" json:"from"`
	ToWalletID   *uuid.UUID `db:"to_wallet_id" json:"to"`
	Amount       float32    `db:"amount" json:"amount"`
}
