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

package runtime

import (
	"github.com/onedomain/lastbackend/pkg/distribution/types"
	"github.com/onedomain/lastbackend/pkg/ingress/envs"
	"github.com/onedomain/lastbackend/pkg/log"
)

func (r Runtime) RouteManage(name string, route *types.RouteManifest) (err error) {

	log.Debugf("route manage: %s", name)

	var status = new(types.RouteStatus)

	defer func() {
		if err = r.config.Sync(); err != nil {
			status.State = types.StateError
			status.Message = err.Error()
			envs.Get().GetState().Routes().SetRouteStatus(name, status)
			return
		}

		if status.State == types.StateDestroy {
			envs.Get().GetState().Routes().DelRoute(name)
			return
		}

		envs.Get().GetState().Routes().SetRouteStatus(name, status)
	}()

	if route.State == types.StateDestroyed {
		status.State = types.StateDestroyed
		envs.Get().GetState().Routes().DelRoute(name)
		return nil
	}

	if route.State == types.StateDestroy {
		status.State = types.StateDestroyed
		envs.Get().GetState().Routes().DelRouteManifests(name)
		return nil
	}

	envs.Get().GetState().Routes().SetRouteManifest(name, route)
	status.State = types.StateProvision

	return nil
}
