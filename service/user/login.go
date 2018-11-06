package user

import (
	"errors"
	"net/http"
	"time"

	"github.com/EdgeSmart/EdgeManager/dao"
	"github.com/EdgeSmart/EdgeManager/token"
	"github.com/gin-gonic/gin"
)

type loginInfo struct {
	Type     string `json:"type"`
	Identity string `json:"identity"`
	Token    string `json:"token"`
}

type response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	MS      int64       `json:"ms"`
}

type loginResponse struct {
	Token string `json:"token"`
}

// Login User login
func Login(ctx *gin.Context) {
	data := loginInfo{}
	err := ctx.BindJSON(&data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"test": "test"})
		return
	}

	uid := ""
	switch data.Type {
	case "test":
		uid, err = loginTest(data)
	}

	tokenStr := ""
	if err == nil {
		tokenObj, err := token.GetInstance("edge")
		if err != nil {
			// todo: error
			return
		}
		tokenStr = tokenObj.Generate(data.Identity)
		err = tokenObj.Set(uid, tokenStr, uid)
		if err != nil {
			// todo: error
			ctx.JSON(500, data)
			return
		}
	}
	responseData := response{
		Status:  0,
		Message: "",
		Data: loginResponse{
			Token: tokenStr,
		},
		MS: time.Now().UnixNano() / 1e6,
	}
	ctx.JSON(200, responseData)
}

func loginTest(data loginInfo) (string, error) {
	db, _ := dao.GetDB("edge")
	row := db.QueryRow("SELECT `uid`,`token`,`ext` FROM `user_auth` WHERE `identity` = ? AND `type` = ? AND `status` = 0", data.Identity, data.Type)

	uid := ""
	token := ""
	var ext string
	err := row.Scan(&uid, &token, &ext)
	if err != nil {
		return "", errors.New("failed")
	}
	if token == data.Token {
		return uid, nil
	}
	return "", errors.New("failed")
}
