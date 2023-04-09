package api

import (
	"main/services/platforms"

	templates "github.com/punk-link/gin-generic-http-templates"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

type StreamingPlatformController struct {
	streamingPlatformService platforms.StreamingPlatformServer
}

func NewStreamingPlatformController(injector *do.Injector) (*StreamingPlatformController, error) {
	streamingPlatformService := do.MustInvoke[platforms.StreamingPlatformServer](injector)

	return &StreamingPlatformController{
		streamingPlatformService: streamingPlatformService,
	}, nil
}

func (t *StreamingPlatformController) RequestUrlSync(ctx *gin.Context) {
	t.streamingPlatformService.PublishPlatforeUrlRequests()

	templates.NoContent(ctx)
}

func (t *StreamingPlatformController) ProcessUrlSyncResults(ctx *gin.Context) {
	t.streamingPlatformService.ProcessPlatforeUrlResults()

	templates.NoContent(ctx)
}
