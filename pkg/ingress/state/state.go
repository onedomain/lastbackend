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

package state

import (
	"github.com/onedomain/lastbackend/pkg/distribution/types"
)

const logLevel = 3

type State struct {
	ingress *IngressState
	routes  *RouteState
}

type IngressState struct {
	Info   types.IngressInfo
	Status types.IngressStatus
}

func (s *State) Ingress() *IngressState {
	return s.ingress
}

func (s *State) Routes() *RouteState {
	return s.routes
}

func New() *State {

	state := State{
		ingress: new(IngressState),
		routes: &RouteState{
			routes: make(map[string]struct {
				status   *types.RouteStatus
				manifest *types.RouteManifest
			}, 0),
			watchers: make(map[chan string]bool, 0),
		},
	}

	return &state
}
