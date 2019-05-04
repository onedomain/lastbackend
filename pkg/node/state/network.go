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

type NetworkState struct {
	lock    sync.RWMutex
	subnets map[string]types.NetworkState
}

func (n *NetworkState) GetSubnets() map[string]types.NetworkState {
	return n.subnets
}

func (n *NetworkState) AddSubnet(cidr string, sn *types.NetworkState) {
	log.V(logLevel).Debugf("Stage: NetworkState: add subnet: %s", cidr)
	n.SetSubnet(cidr, sn)
}

func (n *NetworkState) SetSubnet(cidr string, sn *types.NetworkState) {
	log.V(logLevel).Debugf("Stage: NetworkState: set subnet: %s", cidr)
	n.lock.Lock()
	defer n.lock.Unlock()

	if _, ok := n.subnets[cidr]; ok {
		delete(n.subnets, cidr)
	}

	n.subnets[cidr] = *sn
}

func (n *NetworkState) GetSubnet(cidr string) *types.NetworkState {
	log.V(logLevel).Debugf("Stage: NetworkState: get subnet: %s", cidr)
	n.lock.Lock()
	defer n.lock.Unlock()
	s, ok := n.subnets[cidr]
	if !ok {
		return nil
	}
	return &s
}

func (n *NetworkState) DelSubnet(cidr string) {
	log.V(logLevel).Debugf("Stage: NetworkState: del subnet: %s", cidr)
	n.lock.Lock()
	defer n.lock.Unlock()
	if _, ok := n.subnets[cidr]; ok {
		delete(n.subnets, cidr)
	}
}
