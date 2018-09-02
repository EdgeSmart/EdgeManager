package middleware

import (
	"strings"
	"time"

	"github.com/EdgeSmart/EdgeManager/token"
	"github.com/gin-gonic/gin"
)

// LoginControl 登陆控制
func LoginControl(ctx *gin.Context) {
	if ctx.Writer.Status() == 200 {
		if strings.Index(ctx.Request.RequestURI, "/user/login") != 0 {
			tokenObj, _ := token.GetInstance("edge")
			tokenStr := ctx.DefaultQuery("token", "")
			_, err := tokenObj.Get(tokenStr)
			if err != nil {
				ErrorNeedLogin := struct {
					Code    int    `json:"code"`
					Message string `json:"message"`
					MS      int64  `json:"ms"`
				}{
					Code:    401,
					Message: "Need login",
					MS:      time.Now().UnixNano() / 1e6,
				}
				ctx.JSON(401, ErrorNeedLogin)
				ctx.Abort()
				return
			}
			tokenObj.Active(tokenStr)
		}
	}
	ctx.Next()
	return
}
