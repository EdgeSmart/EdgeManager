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

package mqserver

import (
	"fmt"

	"github.com/EdgeSmart/EdgeManager/service/cluster"
	"github.com/surgemq/surgemq/auth"
)

type edgeAuthenticator struct {
}

var _ auth.Authenticator = (*edgeAuthenticator)(nil)

var (
	edgeAuth = edgeAuthenticator{}
)

func init() {
	fmt.Println("init")
	auth.Register("edge_auth", edgeAuth)
}

func (this edgeAuthenticator) Authenticate(id string, cred interface{}) error {
	return cluster.Auth(id, cred)
}
