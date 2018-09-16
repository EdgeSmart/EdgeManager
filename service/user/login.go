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
	Type     string
	Username string
	Password string
}

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	MS      int64       `json:"ms"`
}

type loginResponse struct {
	Token string `json:"token"`
}

// Login User login
func Login(ctx *gin.Context) {
	var data loginInfo
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
		tokenStr = tokenObj.Generate(data.Username)
		err = tokenObj.Set(uid, tokenStr, uid)
		if err != nil {
			// todo: error
			ctx.JSON(500, data)
			return
		}
	}
	responseData := response{
		Code:    0,
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

	stmt, _ := db.Prepare("SELECT * FROM `user_auth` WHERE `identity` = ? AND `type` = ? AND `status` = ?")
	defer stmt.Close()

	rows := stmt.QueryRow(data.Username, data.Type, 0)

	var uid string
	var token string
	var ext string

	err := rows.Scan(&uid, &token, &ext)
	if err != nil {
		return "", errors.New("failed")
	}
	if token == data.Password {
		return uid, nil
	}
	return "", errors.New("failed")
}
