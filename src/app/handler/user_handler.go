package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"

	"github.com/kujilabo/cocotola-synthesizer-api/src/app/domain"
	"github.com/kujilabo/cocotola-synthesizer-api/src/app/handler/converter"
	"github.com/kujilabo/cocotola-synthesizer-api/src/app/handler/entity"
	handlerhelper "github.com/kujilabo/cocotola-synthesizer-api/src/app/handler/helper"
	"github.com/kujilabo/cocotola-synthesizer-api/src/app/usecase"
	"github.com/kujilabo/cocotola-synthesizer-api/src/lib/ginhelper"
	"github.com/kujilabo/cocotola-synthesizer-api/src/lib/log"
)

type UserHandler interface {
	Synthesize(c *gin.Context)
}

type userHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) UserHandler {
	return &userHandler{userUsecase: userUsecase}
}

func (h *userHandler) Synthesize(c *gin.Context) {
	ctx := c.Request.Context()

	handlerhelper.HandleFunction(c, func() error {
		param := entity.SynthesizeParameter{}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		lang2, err := domain.NewLang2(param.Lang2)
		if err != nil {
			return nil
		}

		result, err := h.userUsecase.Synthesize(ctx, lang2, param.Text)
		if err != nil {
			return xerrors.Errorf("failed to FindSentences. err: %w", err)
		}

		response, err := converter.ToAudioResponse(ctx, result)
		if err != nil {
			return xerrors.Errorf("failed to ToAudioResponse. err: %w", err)
		}

		c.JSON(http.StatusOK, response)
		return nil
	}, h.errorHandle)
}

func (h *userHandler) FindAudioByAudioID(c *gin.Context) {
	ctx := c.Request.Context()

	handlerhelper.HandleFunction(c, func() error {
		audioID, err := ginhelper.GetUintFromPath(c, "audioID")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return nil
		}

		result, err := h.userUsecase.FindAudioByAudioID(ctx, domain.AudioID(audioID))
		if err != nil {
			return xerrors.Errorf("failed to FindSentences. err: %w", err)
		}

		response, err := converter.ToAudioResponse(ctx, result)
		if err != nil {
			return xerrors.Errorf("failed to ToAudioResponse. err: %w", err)
		}

		c.JSON(http.StatusOK, response)
		return nil
	}, h.errorHandle)
}

func (h *userHandler) errorHandle(c *gin.Context, err error) bool {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)

	logger.Errorf("userHandler. err: %v", err)
	return false
}
