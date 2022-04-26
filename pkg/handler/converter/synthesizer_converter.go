package converter

import (
	"context"

	"github.com/kujilabo/cocotola-synthesizer-api/pkg/handler/entity"
	"github.com/kujilabo/cocotola-synthesizer-api/pkg/service"
)

func ToAudioResponse(ctx context.Context, audio service.Audio) (*entity.AudioResponse, error) {
	audioModel := audio.GetAudioModel()
	return &entity.AudioResponse{
		ID:      int(audioModel.GetID()),
		Lang:    audioModel.GetLang().ToLang2().String(),
		Text:    audioModel.GetText(),
		Content: audioModel.GetContent(),
	}, nil
}
