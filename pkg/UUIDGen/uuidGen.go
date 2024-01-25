package uuidgen

import "github.com/gofrs/uuid/v5"

type UUIDGenerator interface {
	New() (uuid.UUID, error)
	FromString(text string) (uuid.UUID, error)
}

type DefaultGenerator struct{}

func (g *DefaultGenerator) New() (uuid.UUID, error) {
	return uuid.NewV4()
}

func (g *DefaultGenerator) FromString(text string) (uuid.UUID, error) {
	return uuid.FromString(text)
}
