package startup

import (
	"main/routers"

	"github.com/gin-gonic/gin"
)

func Configure() *gin.Engine {
	app := gin.Default()
	routers.SetupRouters(app)

	return app
}
