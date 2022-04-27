package usecase

import (
	"context"
	"errors"

	"github.com/kujilabo/cocotola-synthesizer-api/pkg/domain"
	"github.com/kujilabo/cocotola-synthesizer-api/pkg/service"
	"gorm.io/gorm"
)

type UserUsecase interface {
	Synthesize(ctx context.Context, lang2 domain.Lang2, text string) (service.Audio, error)

	FindAudioByAudioID(ctx context.Context, audioID domain.AudioID) (service.Audio, error)
}

type userUsecase struct {
	db                *gorm.DB
	rfFunc            service.RepositoryFactoryFunc
	synthesizerClient service.SynthesizerClient
}

func NewUserUsecase(db *gorm.DB, rfFunc service.RepositoryFactoryFunc, synthesizerClient service.SynthesizerClient) UserUsecase {
	return &userUsecase{
		db:                db,
		rfFunc:            rfFunc,
		synthesizerClient: synthesizerClient,
	}
}

func (u *userUsecase) Synthesize(ctx context.Context, lang2 domain.Lang2, text string) (service.Audio, error) {
	// try to find the audio content from the DB
	rf, err := u.rfFunc(ctx, u.db)
	if err != nil {
		return nil, err
	}

	repo := rf.NewAudioRepository(ctx)
	audioContentFromDB, err := repo.FindByLangAndText(ctx, lang2.ToLang5(), text)
	if err != nil {
		if !errors.Is(err, service.ErrAudioNotFound) {
			return nil, err
		}
	} else {
		return audioContentFromDB, nil
	}

	// synthesize text via the Web API
	audioContent, err := u.synthesizerClient.Synthesize(ctx, lang2.ToLang5(), text)
	if err != nil {
		return nil, err
	}

	audioID, err := repo.AddAudio(ctx, lang2.ToLang5(), text, audioContent)
	if err != nil {
		return nil, err
	}

	audioModel, err := domain.NewAudioModel(uint(audioID), lang2.ToLang5(), text, audioContent)
	if err != nil {
		return nil, err
	}

	return service.NewAudio(audioModel)
}

func (u *userUsecase) FindAudioByAudioID(ctx context.Context, audioID domain.AudioID) (service.Audio, error) {
	rf, err := u.rfFunc(ctx, u.db)
	if err != nil {
		return nil, err
	}

	repo := rf.NewAudioRepository(ctx)
	audio, err := repo.FindAudioByAudioID(ctx, audioID)
	if err != nil {
		return nil, err
	}

	return audio, nil
}
