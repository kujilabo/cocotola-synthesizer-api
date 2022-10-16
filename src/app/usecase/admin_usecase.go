//go:generate mockery --output mock --name AdminUsecase
package usecase

import (
	"github.com/kujilabo/cocotola-synthesizer-api/src/app/service"
)

type AdminUsecase interface {
}

type adminUsecase struct {
	rfFunc service.RepositoryFactoryFunc
}

func NewAdminUsecase(rfFunc service.RepositoryFactoryFunc) AdminUsecase {
	return &adminUsecase{rfFunc: rfFunc}
}
