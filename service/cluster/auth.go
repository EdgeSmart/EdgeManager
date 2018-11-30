package cluster

import (
	"errors"
	"strings"

	"github.com/EdgeSmart/EdgeManager/dao"
)

func Auth(id string, cred interface{}) error {
	cutPos := strings.Index(id, "/")
	if cutPos < 1 {
		return errors.New("Unsupported type")
	}
	idByte := []byte(id)
	authType := string(idByte[0:cutPos])
	switch authType {
	case "manager":
		return nil
	case "gateway":
		return authGateway(string(idByte[cutPos+1:]), cred)
	}
	return errors.New("Unsupported type")
}

func authGateway(id string, cred interface{}) error {
	cutPos := strings.Index(id, "/")
	if cutPos < 1 {
		return errors.New("Auth params error")
	}
	idByte := []byte(id)
	clusterName := string(idByte[0:cutPos])
	machineName := string(idByte[cutPos+1:])
	token := cred.(string)
	if clusterName == "" || machineName == "" || token == "" {
		return errors.New("Auth params error")
	}

	data := dao.GetClusterData(clusterName)
	if data["token"] != token {
		return errors.New("Auth failed")
	}

	return nil
}
