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

package types

// swagger:ignore
type Discovery struct {
	System
	Meta   DiscoveryMeta   `json:"meta"`
	Status DiscoveryStatus `json:"status"`
	Spec   DiscoverySpec   `json:"spec"`
}

type DiscoveryList struct {
	System
	Items []*Discovery
}

type DiscoveryMap struct {
	System
	Items map[string]*Discovery
}

// swagger:ignore
type DiscoveryMeta struct {
	Meta
	SelfLink DiscoverySelfLink `json:"self_link"`
	Node     string            `json:"node"`
}

// swagger:model types_discovery_info
type DiscoveryInfo struct {
	Version      string `json:"version"`
	Hostname     string `json:"hostname"`
	Architecture string `json:"architecture"`

	OSName string `json:"os_name"`
	OSType string `json:"os_type"`

	// RewriteIP - need to set true if you want to use an external ip
	ExternalIP string `json:"external_ip"`
	InternalIP string `json:"internal_ip"`
}

// swagger:model types_discovery_status
type DiscoveryStatus struct {
	IP     string `json:"ip"`
	Port   uint16 `json:"port"`
	Ready  bool   `json:"ready"`
	Online bool   `json:"online"`
}

// swagger:ignore
type DiscoverySpec struct {
}

func (n *Discovery) SelfLink() *DiscoverySelfLink {
	return &n.Meta.SelfLink
}

func NewDiscoveryList() *DiscoveryList {
	dm := new(DiscoveryList)
	dm.Items = make([]*Discovery, 0)
	return dm
}

func NewDiscoveryMap() *DiscoveryMap {
	dm := new(DiscoveryMap)
	dm.Items = make(map[string]*Discovery)
	return dm
}
