//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package request

import (
	"encoding/json"
	"fmt"
	"github.com/onedomain/lastbackend/pkg/api/envs"
	"github.com/onedomain/lastbackend/pkg/distribution/types"
	"gopkg.in/yaml.v2"
	"strings"
)

type RouteManifest struct {
	Meta RouteManifestMeta `json:"meta,omitempty" yaml:"meta,omitempty"`
	Spec RouteManifestSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type RouteManifestMeta struct {
	RuntimeMeta `yaml:",inline"`
}

// swagger:model request_route_create
type RouteManifestSpec struct {
	Port     uint16                         `json:"port" yaml:"port"`
	Type     string                         `json:"type" yaml:"type"`
	Endpoint string                         `json:"endpoint" yaml:"endpoint"`
	Rules    []RouteManifestSpecRulesOption `json:"rules" yaml:"rules"`
}

// swagger:ignore
// swagger:model request_route_remove
type RouteRemoveOptions struct {
	Force bool `json:"force"`
}

// swagger:model request_route_rules
type RouteManifestSpecRulesOption struct {
	Service string `json:"service" yaml:"service"`
	Path    string `json:"path" yaml:"path"`
	Port    int    `json:"port" yaml:"port"`
}

func (r *RouteManifest) FromJson(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *RouteManifest) ToJson() ([]byte, error) {
	return json.Marshal(r)
}

func (r *RouteManifest) FromYaml(data []byte) error {
	return yaml.Unmarshal(data, r)
}

func (r *RouteManifest) ToYaml() ([]byte, error) {
	return yaml.Marshal(r)
}

func (r *RouteManifest) SetRouteMeta(route *types.Route) {

	if route.Meta.Name == types.EmptyString {
		route.Meta.Name = *r.Meta.Name
	}

	if r.Meta.Description != nil {
		route.Meta.Description = *r.Meta.Description
	}

	if r.Meta.Labels != nil {
		route.Meta.Labels = r.Meta.Labels
	}
}

func (r *RouteManifest) SetRouteSpec(route *types.Route, ns *types.Namespace, svc *types.ServiceList) {

	var sl = make(map[string]*types.Service)
	for _, s := range svc.Items {
		sl[s.Meta.Name] = s
	}

	if r.Spec.Endpoint != route.Spec.Endpoint {
		route.Spec.Endpoint = r.Spec.Endpoint
	}

	if r.Spec.Port != route.Spec.Port {
		route.Spec.Port = r.Spec.Port
	}

	if route.Spec.Endpoint == types.EmptyString {
		_, external := envs.Get().GetDomain()
		route.Spec.Endpoint = fmt.Sprintf("%s.%s.%s", strings.ToLower(route.Meta.Name), strings.ToLower(ns.Meta.Name), external)
	}

	route.Spec.Rules = make([]types.RouteRule, 0)
	for _, rs := range r.Spec.Rules {

		if rs.Service == types.EmptyString || rs.Port == 0 {
			continue
		}

		if _, ok := sl[rs.Service]; !ok {
			continue
		}

		route.Spec.Rules = append(route.Spec.Rules, types.RouteRule{
			Upstream: sl[rs.Service].Meta.Endpoint,
			Service:  rs.Service,
			Port:     rs.Port,
			Path:     rs.Path,
		})

	}

	route.Spec.State = types.StateProvision
}
