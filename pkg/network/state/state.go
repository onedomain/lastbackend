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
	subnets   *SubnetState
	endpoints *EndpointState
	resolvers *ResolverState
}

func (s *State) Subnets() *SubnetState {
	return s.subnets
}

func (s *State) Endpoints() *EndpointState {
	return s.endpoints
}

func (s *State) Resolvers () *ResolverState {
	return s.resolvers
}

func New() *State {

	state := State{
		subnets: &SubnetState{
			subnets: make(map[string]types.NetworkState, 0),
		},
		endpoints: &EndpointState{
			endpoints: make(map[string]*types.EndpointState, 0),
		},
		resolvers: &ResolverState{
			resolvers: make(map[string]*types.ResolverManifest, 0),
		},
	}

	return &state
}
