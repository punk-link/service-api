package api

import (
	base "main/controllers"
	"main/services/platforms"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

type StreamingPlatformController struct {
	streamingPlatformService *platforms.StreamingPlatformService
}

func NewStreamingPlatformController(injector *do.Injector) (*StreamingPlatformController, error) {
	streamingPlatformService := do.MustInvoke[*platforms.StreamingPlatformService](injector)

	return &StreamingPlatformController{
		streamingPlatformService: streamingPlatformService,
	}, nil
}

func (t *StreamingPlatformController) RequestUrlSync(ctx *gin.Context) {
	t.streamingPlatformService.PublishPlatforeUrlRequests()

	base.NoContent(ctx)
}

func (t *StreamingPlatformController) ProcessUrlSyncResults(ctx *gin.Context) {
	t.streamingPlatformService.ProcessPlatforeUrlResults()

	base.NoContent(ctx)
}
