package service

import (
	"context"

	"github.com/1boombacks1/testTaskInfotecs/internal/models"
	"github.com/1boombacks1/testTaskInfotecs/internal/repo"
	uuidgen "github.com/1boombacks1/testTaskInfotecs/pkg/UUIDGen"
)

type IService interface {
	CreateWallet(ctx context.Context) (models.Wallet, error)
	Transfer(ctx context.Context, from, to string, amount float32) error
	GetHistoryOperations(ctx context.Context, id string) ([]models.Operations, error)
	GetWalletByID(ctx context.Context, id string) (models.Wallet, error)
}

type Service struct {
	IService
}

type ServicesDependencies struct {
	Repos *repo.Repos
	IDGen uuidgen.UUIDGenerator
}

func NewService(deps ServicesDependencies) *Service {
	return &Service{
		IService: NewWalletService(deps.Repos.Wallet, deps.Repos.Operations, deps.IDGen),
	}
}
