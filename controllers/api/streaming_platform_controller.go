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

func ConstructStreamingPlatformController(injector *do.Injector) (*StreamingPlatformController, error) {
	streamingPlatformService := do.MustInvoke[*platforms.StreamingPlatformService](injector)

	return &StreamingPlatformController{
		streamingPlatformService: streamingPlatformService,
	}, nil
}

func (t *StreamingPlatformController) Sync(ctx *gin.Context) {
	t.streamingPlatformService.PublishPlatforeUrlRequests()

	base.NoContent(ctx)
}
