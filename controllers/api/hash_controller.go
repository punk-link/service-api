package api

import (
	base "main/controllers"
	"main/services/common"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

type HashController struct {
	coder *common.HashCoder
}

func NewHashController(injector *do.Injector) (*HashController, error) {
	coder := do.MustInvoke[*common.HashCoder](injector)

	return &HashController{
		coder: coder,
	}, nil
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
