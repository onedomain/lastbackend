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

package endpoint

import (
	"context"

	"github.com/onedomain/lastbackend/pkg/discovery/envs"
	"github.com/onedomain/lastbackend/pkg/distribution"
	"github.com/onedomain/lastbackend/pkg/distribution/types"
	"github.com/onedomain/lastbackend/pkg/log"
)

const (
	logLevel  = 3
	logPrefix = "runtime:endpoint"
)

func Watch(ctx context.Context) {

	log.V(logLevel).Debugf("%s:restore:> watch change endpoint start", logPrefix)

	var (
		em    = distribution.NewEndpointModel(ctx, envs.Get().GetStorage())
		cache = envs.Get().GetCache().Endpoint()
		event = make(chan types.EndpointEvent)
	)

	go func() {
		for {
			select {
			case e := <-event:
				{

					if e.Data == nil {
						continue
					}

					endpoint := e.Data

					switch e.Action {
					case types.EventActionCreate:
						fallthrough
					case types.EventActionUpdate:
						cache.Del(endpoint.Spec.Domain)
						envs.Get().GetCache().Endpoint().Set(endpoint.Spec.Domain, []string{endpoint.Spec.IP})
						continue
					case types.EventActionDelete:
						cache.Del(endpoint.Spec.Domain)
						continue
					}

				}
			}
		}
	}()

	go em.Watch(event, nil)
}
