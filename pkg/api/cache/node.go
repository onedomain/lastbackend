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
	"fmt"
	"sync"

	"github.com/onedomain/lastbackend/pkg/distribution/types"
	"github.com/onedomain/lastbackend/pkg/log"
)

const logCacheNode = "api:cache:node"

type CacheNodeManifest struct {
	lock      sync.RWMutex
	nodes     map[string]*types.Node
	ingress   map[string]*types.Ingress
	exporter  map[string]*types.Exporter
	discovery map[string]*types.Discovery
	configs   map[string]*types.ConfigManifest
	manifests map[string]*types.NodeManifest
}

func (c *CacheNodeManifest) checkNode(node string) {
	if _, ok := c.manifests[node]; !ok {
		c.manifests[node] = new(types.NodeManifest)
	}
}

func (c *CacheNodeManifest) SetPodManifest(node, pod string, s *types.PodManifest) {
	log.Infof("%s:PodManifestSet:> %s, %s, %#v", logCacheNode, node, pod, s)
	c.lock.Lock()
	defer c.lock.Unlock()

	c.checkNode(node)

	if c.manifests[node].Pods == nil {
		sp := c.manifests[node]
		sp.Pods = make(map[string]*types.PodManifest, 0)
	}

	c.manifests[node].Pods[pod] = s
}

func (c *CacheNodeManifest) DelPodManifest(node, pod string) {
	log.Infof("%s:PodManifestDel:> %s, %s", logCacheNode, node, pod)
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, ok := c.manifests[node]; !ok {
		return
	}

	delete(c.manifests[node].Pods, pod)
}

func (c *CacheNodeManifest) SetVolumeManifest(node, volume string, s *types.VolumeManifest) {

	log.Infof("%s:SetVolumeManifest:> %s, %s", logCacheNode, node, volume)

	c.lock.Lock()
	defer c.lock.Unlock()

	c.checkNode(node)

	if c.manifests[node].Volumes == nil {
		sp := c.manifests[node]
		sp.Volumes = make(map[string]*types.VolumeManifest, 0)
	}

	c.manifests[node].Volumes[volume] = s
}

func (c *CacheNodeManifest) DelVolumeManifest(node, volume string) {

	log.Infof("%s:DelVolumeManifest:> %s, %s", logCacheNode, node, volume)

	c.lock.Lock()
	defer c.lock.Unlock()

	if _, ok := c.manifests[node]; !ok {
		return
	}

	delete(c.manifests[node].Volumes, volume)
}

func (c *CacheNodeManifest) SetSubnetManifest(cidr string, s *types.SubnetManifest) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for n := range c.manifests {

		if _, ok := c.manifests[n].Network[cidr]; !ok {
			c.manifests[n].Network = make(map[string]*types.SubnetManifest)
		}

		c.manifests[n].Network[cidr] = s
	}
}

func (c *CacheNodeManifest) SetSecretManifest(name string, s *types.SecretManifest) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for n := range c.manifests {

		if _, ok := c.manifests[n].Secrets[name]; !ok {
			c.manifests[n].Secrets = make(map[string]*types.SecretManifest)
		}

		c.manifests[n].Secrets[name] = s
	}
}

func (c *CacheNodeManifest) SetConfigManifest(name string, s *types.ConfigManifest) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.configs[name] = s
	for n := range c.manifests {
		if _, ok := c.manifests[n].Configs[name]; !ok {
			c.manifests[n].Configs = make(map[string]*types.ConfigManifest)
		}

		c.manifests[n].Configs[name] = s
	}
}

func (c *CacheNodeManifest) SetEndpointManifest(addr string, s *types.EndpointManifest) {
	c.lock.Lock()
	defer c.lock.Unlock()

	log.Debugf("%s set endpoint manifest: %s > %s", logCacheNode, addr, s.IP)

	for _, n := range c.manifests {
		if n.Endpoints == nil {
			n.Endpoints = make(map[string]*types.EndpointManifest, 0)
		}
		n.Endpoints[addr] = s
	}
}

func (c *CacheNodeManifest) SetIngress(ingress *types.Ingress) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.ingress[ingress.SelfLink().String()] = ingress
}

func (c *CacheNodeManifest) DelIngress(selflink string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.ingress, selflink)
}

func (c *CacheNodeManifest) SetDiscovery(discovery *types.Discovery) {
	c.lock.Lock()
	defer c.lock.Unlock()

	dvc, ok := c.discovery[discovery.SelfLink().String()]

	if !ok {
		c.discovery[discovery.SelfLink().String()] = discovery
		c.SetResolvers()
		return
	}

	var update = false
	switch true {
	case dvc.Status.IP != discovery.Status.IP:
		update = true
		break
	case dvc.Status.Port != discovery.Status.Port:
		update = true
		break
	case dvc.Status.Ready != discovery.Status.Ready:
		update = true
		break
	}
	if update {
		c.discovery[discovery.SelfLink().String()] = discovery
		c.SetResolvers()
	}
	return
}

func (c *CacheNodeManifest) DelDiscovery(selflink string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.discovery, selflink)

	resolvers := make(map[string]*types.ResolverManifest, 0)

	for _, d := range c.discovery {
		if d.Status.Ready {
			resolvers[d.Status.IP] = &types.ResolverManifest{
				IP:   d.Status.IP,
				Port: d.Status.Port,
			}
		}
	}

	for _, n := range c.manifests {
		n.Resolvers = resolvers
	}
}

func (c *CacheNodeManifest) SetExporter(exporter *types.Exporter) {
	c.lock.Lock()
	defer c.lock.Unlock()

	dvc, ok := c.exporter[exporter.SelfLink().String()]

	if !ok {
		c.exporter[exporter.SelfLink().String()] = exporter
		c.SetExporterEndpoint()
		return
	}

	var update = false
	switch true {
	case dvc.Status.Listener.IP != exporter.Status.Listener.IP:
		update = true
		break
	case dvc.Status.Listener.Port != exporter.Status.Listener.Port:
		update = true
		break
	case dvc.Status.Ready != exporter.Status.Ready:
		update = true
		break
	}
	if update {
		c.exporter[exporter.SelfLink().String()] = exporter
		c.SetExporterEndpoint()
	}
	return
}

func (c *CacheNodeManifest) DelExporter(selflink string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.exporter, selflink)

	for _, n := range c.manifests {
		n.Exporter = nil
	}

	for _, d := range c.exporter {
		if d.Status.Ready {

			exporter := &types.ExporterManifest{
				Endpoint: fmt.Sprintf("%s:%d", d.Status.Listener.IP, d.Status.Listener.Port),
			}

			for _, n := range c.manifests {
				n.Exporter = exporter
			}

			break
		}
	}
}

func (c *CacheNodeManifest) SetExporterEndpoint() {

	for _, n := range c.manifests {
		n.Exporter = nil
	}

	for _, d := range c.exporter {
		if d.Status.Ready {

			exporter := &types.ExporterManifest{
				Endpoint: fmt.Sprintf("%s:%d", d.Status.Listener.IP, d.Status.Listener.Port),
			}

			for _, n := range c.manifests {
				n.Exporter = exporter
			}

			break
		}
	}
}

func (c *CacheNodeManifest) GetExporterEndpoint() *types.ExporterManifest {

	c.lock.Lock()
	defer c.lock.Unlock()

	exporter := new(types.ExporterManifest)

	for _, d := range c.exporter {
		if d.Status.Ready {
			exporter.Endpoint = fmt.Sprintf("%s:%d", d.Status.Listener.IP, d.Status.Listener.Port)
		}
	}

	return exporter
}

func (c *CacheNodeManifest) SetResolvers() {
	resolvers := make(map[string]*types.ResolverManifest, 0)

	for _, d := range c.discovery {
		if d.Status.Ready {
			resolvers[d.Status.IP] = &types.ResolverManifest{
				IP:   d.Status.IP,
				Port: d.Status.Port,
			}
		}
	}

	for _, n := range c.manifests {
		n.Resolvers = resolvers
	}
}

func (c *CacheNodeManifest) GetResolvers() map[string]*types.ResolverManifest {

	resolvers := make(map[string]*types.ResolverManifest, 0)

	for _, d := range c.discovery {
		if d.Status.Ready {
			resolvers[d.Status.IP] = &types.ResolverManifest{
				IP:   d.Status.IP,
				Port: d.Status.Port,
			}
		}
	}

	return resolvers
}

func (c *CacheNodeManifest) GetConfigs() map[string]*types.ConfigManifest {
	return c.configs
}

func (c *CacheNodeManifest) SetNode(node *types.Node) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.nodes[node.SelfLink().String()] = node
}

func (c *CacheNodeManifest) DelNode(node *types.Node) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.nodes, node.SelfLink().String())
	delete(c.manifests, node.SelfLink().String())
}

func (c *CacheNodeManifest) Get(node string) *types.NodeManifest {
	c.lock.Lock()
	defer c.lock.Unlock()
	if s, ok := c.manifests[node]; !ok {
		return nil
	} else {
		return s
	}
}

func (c *CacheNodeManifest) Flush(node string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.manifests[node] = new(types.NodeManifest)
}

func (c *CacheNodeManifest) Clear(node string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.manifests, node)
}

func NewCacheNodeManifest() *CacheNodeManifest {
	c := new(CacheNodeManifest)
	c.exporter = make(map[string]*types.Exporter, 0)
	c.manifests = make(map[string]*types.NodeManifest, 0)
	c.ingress = make(map[string]*types.Ingress, 0)
	c.discovery = make(map[string]*types.Discovery, 0)
	c.configs = make(map[string]*types.ConfigManifest, 0)
	return c
}
