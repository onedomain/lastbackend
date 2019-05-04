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

package cache

import (
	"github.com/onedomain/lastbackend/pkg/log"
	"sync"

	"github.com/onedomain/lastbackend/pkg/distribution/types"
)

const logCacheDiscovery = "api:cache:discovery"

type CacheDiscoveryManifest struct {
	lock      sync.RWMutex
	manifests map[string]*types.DiscoveryManifest
}

func (c *CacheDiscoveryManifest) SetSubnetManifest(cidr string, s *types.SubnetManifest) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for n := range c.manifests {

		if _, ok := c.manifests[n].Network[cidr]; !ok {
			c.manifests[n].Network = make(map[string]*types.SubnetManifest)
		}

		c.manifests[n].Network[cidr] = s
	}
}

func (c *CacheDiscoveryManifest) Get(discovery string) *types.DiscoveryManifest {
	c.lock.Lock()
	defer c.lock.Unlock()
	if s, ok := c.manifests[discovery]; !ok {
		return nil
	} else {
		return s
	}
}

func (c *CacheDiscoveryManifest) Flush(discovery string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.manifests[discovery] = new(types.DiscoveryManifest)
}

func (c *CacheDiscoveryManifest) Clear(discovery string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	log.Debugf("clear cache for discovery: %s", discovery)
	delete(c.manifests, discovery)
}

func NewCacheDiscoveryManifest() *CacheDiscoveryManifest {
	c := new(CacheDiscoveryManifest)
	c.manifests = make(map[string]*types.DiscoveryManifest, 0)
	return c
}
