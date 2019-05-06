// Copyright (c) 2014 The SurgeMQ Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mqtt

import (
	"errors"
	"strings"

	"github.com/EdgeSmart/EdgeManager/dao"
	"github.com/surgemq/surgemq/auth"
)

type (
	edgeAuth struct{}
)

var (
	_ auth.Authenticator = (*edgeAuth)(nil)
)

func init() {
	auth.Register("edge_auth", &edgeAuth{})
}

// Authenticate Authenticate
func (e *edgeAuth) Authenticate(id string, cred interface{}) error {
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
		return e.authGateway(string(idByte[cutPos+1:]), cred)
	}
	return errors.New("Unsupported type")
}

// authGateway authGateway
func (e *edgeAuth) authGateway(id string, cred interface{}) error {
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
