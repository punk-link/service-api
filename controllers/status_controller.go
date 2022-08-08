package controllers

import (
	"main/data"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func CheckHealth(ctx *gin.Context) {
	sqlDb, err := data.DB.DB()
	if err != nil {
		log.Error().Err(err).Msgf("Postgres initialization failed: %v", err.Error())
	}

	err = sqlDb.Ping()
	if err != nil {
		log.Error().Err(err).Msgf("Can't reach any database instances: %v", err.Error())
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}
