-- +goose Up
CREATE TABLE wallets(
    id uuid PRIMARY KEY,
    balance FLOAT not null default 100.0,
    created_at TIMESTAMP not null DEFAULT now(),
    updated_at TIMESTAMP not null DEFAULT now()
);

CREATE TABLE operations(
    id serial PRIMARY KEY,
    from_wallet_id uuid not null,
    to_wallet_id uuid not null,
    amount FLOAT not null,
    created_at TIMESTAMP not null DEFAULT now(),
    FOREIGN KEY (from_wallet_id) REFERENCES wallets (id) ON DELETE CASCADE,
    FOREIGN KEY (to_wallet_id) REFERENCES wallets (id) ON DELETE CASCADE
);

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS operations;
drop TABLE if EXISTS wallets;
-- +goose StatementEnd
