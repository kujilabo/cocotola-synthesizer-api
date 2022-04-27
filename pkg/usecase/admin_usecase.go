package usecase

import (
	"github.com/kujilabo/cocotola-synthesizer-api/pkg/service"
)

type AdminUsecase interface {
}

type adminUsecase struct {
	rf service.RepositoryFactory
}

func NewAdminUsecase(rf service.RepositoryFactory) AdminUsecase {
	return &adminUsecase{rf: rf}
}
