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
type Ingress struct {
	System
	Meta   IngressMeta   `json:"meta"`
	Status IngressStatus `json:"status"`
	Spec   IngressSpec   `json:"spec"`
}

type IngressList struct {
	System
	Items []*Ingress
}

type IngressMap struct {
	System
	Items map[string]*Ingress
}

// swagger:ignore
type IngressMeta struct {
	Meta
	SelfLink IngressSelfLink `json:"self_link"`
	Node     string          `json:"node"`
}

type IngressInfo struct {
	Type         string `json:"type"`
	Version      string `json:"version"`
	Hostname     string `json:"hostname"`
	Architecture string `json:"architecture"`

	OSName string `json:"os_name"`
	OSType string `json:"os_type"`

	// RewriteIP - need to set true if you want to use an external ip
	ExternalIP string `json:"external_ip"`
	InternalIP string `json:"internal_ip"`
}

// swagger:model types_ingress_status
type IngressStatus struct {
	Ready  bool `json:"ready"`
	Online bool `json:"online"`
}

type IngressSpec struct {
}

// swagger:ignore
func (n *Ingress) SelfLink() *IngressSelfLink {
	return &n.Meta.SelfLink
}

func NewIngressList() *IngressList {
	dm := new(IngressList)
	dm.Items = make([]*Ingress, 0)
	return dm
}

func NewIngressMap() *IngressMap {
	dm := new(IngressMap)
	dm.Items = make(map[string]*Ingress)
	return dm
}
