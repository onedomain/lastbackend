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

package views

import (
	"encoding/json"

	"github.com/onedomain/lastbackend/pkg/distribution/types"
)

type RouteView struct{}

func (rv *RouteView) New(obj *types.Route) *Route {
	r := Route{}
	r.Meta = r.ToMeta(obj.Meta)
	r.Spec = r.ToSpec(obj.Spec)
	r.Status = r.ToStatus(obj.Status)
	return &r
}

func (p *Route) ToJson() ([]byte, error) {
	return json.Marshal(p)
}

func (r *Route) ToMeta(obj types.RouteMeta) RouteMeta {
	meta := RouteMeta{}
	meta.Name = obj.Name
	meta.Namespace = obj.Namespace
	meta.SelfLink = obj.SelfLink.String()
	meta.Updated = obj.Updated
	meta.Created = obj.Created

	return meta
}

func (r *Route) ToSpec(obj types.RouteSpec) RouteSpec {
	spec := RouteSpec{}
	spec.Domain = obj.Endpoint
	spec.Port = obj.Port
	for _, rule := range obj.Rules {
		spec.Rules = append(spec.Rules, &RouteRule{
			Service:  rule.Service,
			Path:     rule.Path,
			Port:     rule.Port,
			Endpoint: rule.Upstream,
		})
	}
	return spec
}

func (r *Route) ToStatus(obj types.RouteStatus) RouteStatus {
	state := RouteStatus{}
	state.State = obj.State
	state.Message = obj.Message
	return state
}

func (rv RouteView) NewList(obj *types.RouteList) *RouteList {
	if obj == nil {
		return nil
	}

	n := make(RouteList, 0)
	for _, v := range obj.Items {
		n = append(n, rv.New(v))
	}
	return &n
}

func (n *RouteList) ToJson() ([]byte, error) {
	if n == nil {
		n = &RouteList{}
	}
	return json.Marshal(n)
}
