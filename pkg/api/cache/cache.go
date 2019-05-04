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

import "context"

const logLevel = 3

type Cache struct {
	node      *CacheNodeManifest
	ingress   *CacheIngressManifest
	discovery *CacheDiscoveryManifest
	exporter  *CacheExporterManifest
}

type Cleaner func(ctx context.Context) error

func (c *Cache) Node() *CacheNodeManifest {
	return c.node
}

func (c *Cache) Ingress() *CacheIngressManifest {
	return c.ingress
}

func (c *Cache) Discovery() *CacheDiscoveryManifest {
	return c.discovery
}

func (c *Cache) Exporter() *CacheExporterManifest {
	return c.exporter
}

func NewCache() *Cache {
	c := new(Cache)
	c.node = NewCacheNodeManifest()
	c.ingress = NewCacheIngressManifest()
	c.discovery = NewCacheDiscoveryManifest()
	c.exporter = NewCacheExporterManifest()
	return c
}
