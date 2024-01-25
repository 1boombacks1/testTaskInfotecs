package service

import "fmt"

var (
	ErrCannotCreateWallet = fmt.Errorf("cannot create wallet")
	ErrWalletNotFound     = fmt.Errorf("wallet not found")
	ErrCannotGetWallet    = fmt.Errorf("cannot get account")

	ErrFromWalletNotFound        = fmt.Errorf("outgoing wallet not found")
	ErrToWalletNotFound          = fmt.Errorf("recipient's wallet not found")
	ErrNotEnoughBalanceToTranser = fmt.Errorf("not enough balance to transfer")
	ErrTrasfer                   = fmt.Errorf("transfer error")

	ErrInvalidID = fmt.Errorf("invalid input id. Check it for correctness")
)
