package converter

import (
	"context"

	"github.com/kujilabo/cocotola-synthesizer-api/src/app/controller/entity"
	"github.com/kujilabo/cocotola-synthesizer-api/src/app/service"
)

func ToAudioResponse(ctx context.Context, audio service.Audio) (*entity.AudioResponse, error) {
	audioModel := audio.GetAudioModel()
	return &entity.AudioResponse{
		ID:      int(audioModel.GetID()),
		Lang2:   audioModel.GetLang5().ToLang2().String(),
		Text:    audioModel.GetText(),
		Content: audioModel.GetContent(),
	}, nil
}
