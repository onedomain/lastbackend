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

package controller

import (
	"context"
	"github.com/onedomain/lastbackend/pkg/api/types/v1"
	"github.com/onedomain/lastbackend/pkg/api/types/v1/request"
	"github.com/onedomain/lastbackend/pkg/distribution/types"
	"github.com/onedomain/lastbackend/pkg/exporter/envs"
	"github.com/onedomain/lastbackend/pkg/exporter/runtime"
	"github.com/onedomain/lastbackend/pkg/log"
	"sync"
	"time"
)

const (
	logPrefix = "controller:>"
	logLevel  = 3
)

type Controller struct {
	runtime *runtime.Runtime
	cache   struct {
		lock      sync.RWMutex
		resources types.ExporterStatus
	}
}

func New(r *runtime.Runtime) *Controller {
	var c = new(Controller)
	c.runtime = r
	return c
}

func (c *Controller) Connect(ctx context.Context) error {

	log.V(logLevel).Debugf("%s:connect:> connect init", logPrefix)

	opts := v1.Request().Exporter().ExporterConnectOptions()
	opts.Info = envs.Get().GetState().Exporter().Info
	opts.Status = envs.Get().GetState().Exporter().Status

	for {
		err := envs.Get().GetClient().Connect(ctx, opts)
		if err == nil {
			log.Debugf("%s connected", logPrefix)
			return nil
		}

		log.Errorf("connect err: %s", err.Error())
		time.Sleep(3 * time.Second)
	}

}

func (c *Controller) Sync() error {

	log.Debugf("Start exporter sync")

	ticker := time.NewTicker(time.Second * 5)

	for range ticker.C {
		opts := new(request.ExporterStatusOptions)
		status := envs.Get().GetState().Exporter().Status
		opts.Listener.IP = status.Listener.IP
		opts.Listener.Port = status.Listener.Port
		opts.Http.IP = status.Http.IP
		opts.Http.Port = status.Http.Port
		opts.Ready = status.Ready

		c.cache.lock.Lock()

		spec, err := envs.Get().GetClient().SetStatus(context.Background(), opts)
		if err != nil {
			log.Errorf("exporter:exporter:dispatch err: %s", err.Error())
		}

		c.cache.lock.Unlock()
		if spec != nil {

		} else {
			log.Debug("received spec is nil, skip apply changes")
		}
	}

	return nil
}
