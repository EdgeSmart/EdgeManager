package user

import (
	"github.com/gin-gonic/gin"
)

func Test(ctx *gin.Context) {
	data := map[string]string{
		"data": "sadfsadfsadf",
		"aaaa": "ewrewrew",
	}
	ctx.JSON(200, data)
}
