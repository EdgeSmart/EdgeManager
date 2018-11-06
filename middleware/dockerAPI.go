package middleware

import (
	"github.com/EdgeSmart/EdgeManager/mqserver"
	"github.com/gin-gonic/gin"
)

// DockerAPI DockerAPI
func DockerAPI(ctx *gin.Context) {
	mqserver.HandleClientV2(ctx)
	// ctx.Next()
	return
}
