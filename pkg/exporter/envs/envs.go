//
// KULADO INC. CONFIDENTIAL
// __________________
//
// [2014] - [2019] KULADO INC.
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of KULADO INC. and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to KULADO INC.
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from KULADO INC..
//

package envs

import (
	"github.com/onedomain/lastbackend/pkg/api/client/types"
	"github.com/onedomain/lastbackend/pkg/exporter/logger"
	"github.com/onedomain/lastbackend/pkg/exporter/state"
)

var _env Env

type Env struct {
	state       *state.State
	logger      *logger.Logger
	client      types.ExporterClientV1
	accessToken string
}

func Get() *Env {
	return &_env
}

func (c *Env) SetState(state *state.State) {
	c.state = state
}

func (c *Env) GetState() *state.State {
	return c.state
}

func (c *Env) SetLogger(logger *logger.Logger) {
	c.logger = logger
}

func (c *Env) GetLogger() *logger.Logger {
	return c.logger
}

func (c *Env) SetClient(client types.ExporterClientV1) {
	c.client = client
}

func (c *Env) GetClient() types.ExporterClientV1 {
	return c.client
}

func (c *Env) SetAccessToken(token string) {
	c.accessToken = token
}

func (c *Env) GetAccessToken() string {
	return c.accessToken
}
