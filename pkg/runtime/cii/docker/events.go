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

package docker

import (
	"context"
	d "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/onedomain/lastbackend/pkg/distribution/types"
	"github.com/onedomain/lastbackend/pkg/log"
)

func (r *Runtime) Subscribe(ctx context.Context) (chan *types.Image, error) {

	log.V(logLevel).Debug("Create new event listener subscribe")
	var cs = make(chan *types.Image)

	go func() {

		if _, err := r.client.Ping(ctx); err != nil {
			log.Errorf("Can not ping docker client")
			return
		}

		es, errr := r.client.Events(ctx, d.EventsOptions{})
		for {
			select {
			case e := <-es:

				if e.Type != events.ImageEventType {
					continue
				}

				log.V(logLevel).Debugf("Image %s", e.ID)

				if e.Action == types.StateDestroy {
					c := new(types.Image)
					c.Meta.ID = e.ID
					c.Status.State = types.StateDestroyed
					cs <- c
					break
				}

				c, err := r.Inspect(ctx, e.ID)
				if err != nil {
					log.Errorf("Container inspect err: %s", err.Error())
					continue
				}
				if c == nil {
					log.Errorf("Container: container not found")
					break
				}
				break

			case err := <-errr:
				log.Errorf("Event listening error: %s", err)
			}
		}
	}()

	return cs, nil

}
