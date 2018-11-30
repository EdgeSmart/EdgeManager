package cluster

import (
	"fmt"
	"net/http"

	"github.com/EdgeSmart/EdgeManager/dao"
	"github.com/gin-gonic/gin"
)

func List(ctx *gin.Context) {

	type reqParam struct {
		Token string `json:"token"`
	}

	type resParam struct {
		CID   string `json:"cid"`
		CKey  string `json:"ckey"`
		UID   string `json:"uid"`
		Name  string `json:"name"`
		Token string `json:"token"`
	}

	type retStruct struct {
		Status  int
		Message string
		Data    interface{}
	}

	db, _ := dao.GetDB("edge")
	rows, err := db.Query("SELECT cid,ckey,uid,name,token FROM `cluster` WHERE `status` = 1")
	if err != nil {
		fmt.Println("cluster_list_err", err)
		return
	}
	cid := ""
	ckey := ""
	uid := ""
	name := ""
	token := ""
	reqData := []resParam{}
	for rows.Next() {
		err := rows.Scan(&cid, &ckey, &uid, &name, &token)
		if err != nil {
			fmt.Println("cluster_list_err", err)
			return
		}
		data := resParam{
			CID:   cid,
			CKey:  ckey,
			UID:   uid,
			Name:  name,
			Token: token,
		}
		reqData = append(reqData, data)
	}
	ret := retStruct{
		Status:  0,
		Message: "",
		Data:    reqData,
	}
	ctx.JSON(http.StatusOK, ret)

}
