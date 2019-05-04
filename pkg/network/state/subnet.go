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

const logSubnetPrefix = "state:subnet:>"

type SubnetState struct {
	lock    sync.RWMutex
	subnets map[string]types.NetworkState
}

func (n *SubnetState) GetSubnets() map[string]types.NetworkState {
	return n.subnets
}

func (n *SubnetState) AddSubnet(cidr string, sn *types.NetworkState) {
	log.V(logLevel).Debugf("%s add subnet: %s", logSubnetPrefix, cidr)
	n.SetSubnet(cidr, sn)
}

func (n *SubnetState) SetSubnet(cidr string, sn *types.NetworkState) {
	log.V(logLevel).Debugf("%s set subnet: %s", logSubnetPrefix, cidr)
	n.lock.Lock()
	defer n.lock.Unlock()

	if _, ok := n.subnets[cidr]; ok {
		delete(n.subnets, cidr)
	}

	n.subnets[cidr] = *sn
}

func (n *SubnetState) GetSubnet(cidr string) *types.NetworkState {
	log.V(logLevel).Debugf("%s get subnet: %s", logSubnetPrefix, cidr)
	n.lock.Lock()
	defer n.lock.Unlock()
	s, ok := n.subnets[cidr]
	if !ok {
		return nil
	}
	return &s
}

func (n *SubnetState) DelSubnet(cidr string) {
	log.V(logLevel).Debugf("%s del subnet: %s", logSubnetPrefix, cidr)
	n.lock.Lock()
	defer n.lock.Unlock()
	if _, ok := n.subnets[cidr]; ok {
		delete(n.subnets, cidr)
	}
}
