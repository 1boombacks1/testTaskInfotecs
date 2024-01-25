package service

import (
	"context"
	"errors"

	"github.com/1boombacks1/testTaskInfotecs/internal/models"
	"github.com/1boombacks1/testTaskInfotecs/internal/repo"
	"github.com/1boombacks1/testTaskInfotecs/internal/repo/repoerrs"
	uuidgen "github.com/1boombacks1/testTaskInfotecs/pkg/UUIDGen"
)

type WalletService struct {
	walletRepo     repo.Wallet
	operationsRepo repo.Operations
	idGen          uuidgen.UUIDGenerator
}

func NewWalletService(wr repo.Wallet, or repo.Operations, idGenerator uuidgen.UUIDGenerator) *WalletService {
	return &WalletService{
		walletRepo:     wr,
		operationsRepo: or,
		idGen:          idGenerator,
	}
}

func (s *WalletService) CreateWallet(ctx context.Context) (models.Wallet, error) {
	id, err := s.idGen.New()
	if err != nil {
		return models.Wallet{}, ErrCannotCreateWallet
	}
	err = s.walletRepo.Create(ctx, id)
	if err != nil {
		return models.Wallet{}, ErrCannotCreateWallet
	}

	return models.Wallet{ID: id, Balance: 100.0}, nil
}

func (s *WalletService) Transfer(ctx context.Context, from, to string, amount float32) error {
	fromParsed, err := s.idGen.FromString(from)
	if err != nil {
		return ErrInvalidID
	}

	_, err = s.walletRepo.GetByID(ctx, fromParsed)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return ErrFromWalletNotFound
		}
		return ErrTrasfer
	}

	toParsed, err := s.idGen.FromString(to)
	if err != nil {
		return ErrInvalidID
	}

	if _, err = s.walletRepo.GetByID(ctx, toParsed); err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return ErrToWalletNotFound
		}
		return ErrTrasfer
	}

	if err = s.walletRepo.Transfer(ctx, fromParsed, toParsed, amount); err != nil {
		if errors.Is(err, repoerrs.ErrNotEnoughBalance) {
			return ErrNotEnoughBalanceToTranser
		}
		return ErrTrasfer
	}

	if _, err = s.operationsRepo.Create(ctx, fromParsed, toParsed, amount); err != nil {
		return err
	}

	return nil
}

func (s *WalletService) GetHistoryOperations(ctx context.Context, id string) ([]models.Operations, error) {
	parsedID, err := s.idGen.FromString(id)
	if err != nil {
		return nil, ErrInvalidID
	}
	operations, err := s.operationsRepo.GetAllByWalletID(ctx, parsedID)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return nil, ErrWalletNotFound
		}
		return nil, err
	}
	return operations, nil
}

func (s *WalletService) GetWalletByID(ctx context.Context, id string) (models.Wallet, error) {
	parsedID, err := s.idGen.FromString(id)
	if err != nil {
		return models.Wallet{}, ErrInvalidID
	}
	wallet, err := s.walletRepo.GetByID(ctx, parsedID)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return models.Wallet{}, ErrWalletNotFound
		}
		return models.Wallet{}, err
	}
	return wallet, nil
}
