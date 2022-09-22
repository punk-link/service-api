package api

import (
	base "main/controllers"
	"main/services/common"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HashController struct {
	coder *common.HashCoder
}

func ConstructHashController(coder *common.HashCoder) *HashController {
	return &HashController{
		coder: coder,
	}
}

func (t *HashController) Decode(ctx *gin.Context) {
	hash := ctx.Param("target")

	id := t.coder.Decode(hash)
	base.Ok(ctx, id)
}

func (t *HashController) Encode(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("target"))

	hash := t.coder.Encode(id)
	base.OkOrBadRequest(ctx, hash, err)
}
