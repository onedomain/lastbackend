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
	"github.com/onedomain/lastbackend/pkg/log"
	"sync"
)

const logResolversPrefix = "state:resolvers:>"

type ResolverState struct {
	lock      sync.RWMutex
	resolvers map[string]*types.ResolverManifest
}

func (n *ResolverState) GetResolvers() map[string]*types.ResolverManifest {
	return n.resolvers
}

func (n *ResolverState) AddResolver(cidr string, sn *types.ResolverManifest) {
	log.V(logLevel).Debugf("%s add resolver: %s", logResolversPrefix, cidr)
	n.SetResolver(cidr, sn)
}

func (n *ResolverState) SetResolver(cidr string, sn *types.ResolverManifest) {
	log.V(logLevel).Debugf("%s set resolver: %s", logResolversPrefix, cidr)
	n.lock.Lock()
	defer n.lock.Unlock()

	if _, ok := n.resolvers[cidr]; ok {
		delete(n.resolvers, cidr)
	}

	n.resolvers[cidr] = sn
}

func (n *ResolverState) GetResolver(cidr string) *types.ResolverManifest {
	log.V(logLevel).Debugf("%s get resolver: %s", logResolversPrefix, cidr)
	n.lock.Lock()
	defer n.lock.Unlock()
	s, ok := n.resolvers[cidr]
	if !ok {
		return nil
	}
	return s
}

func (n *ResolverState) DelResolver(cidr string) {
	log.V(logLevel).Debugf("%s del resolver: %s", logResolversPrefix, cidr)
	n.lock.Lock()
	defer n.lock.Unlock()
	if _, ok := n.resolvers[cidr]; ok {
		delete(n.resolvers, cidr)
	}
}
