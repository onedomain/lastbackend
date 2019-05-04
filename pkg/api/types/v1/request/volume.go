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

package request

import (
	"encoding/json"
	"github.com/onedomain/lastbackend/pkg/distribution/types"
	"github.com/onedomain/lastbackend/pkg/util/resource"
	"gopkg.in/yaml.v2"
	"time"
)

type VolumeManifest struct{
	Meta VolumeManifestMeta `json:"meta,omitempty" yaml:"meta,omitempty"`
	Spec VolumeManifestSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type VolumeManifestMeta struct {
	RuntimeMeta `yaml:",inline"`
}

type VolumeManifestSpec struct {
	// Template volume types
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
	// Template volume selector
	Selector ManifestSpecSelector `json:"selector,omitempty" yaml:"selector,omitempty"`
	//  Volume Resources
	Capacity VolumeManifestSpecCapacity `json:"capacity,omitempty" yaml:"capacity,omitempty"`
	// Volume access mode
	AccessMode string `json:"access_mode,omitempty" yaml:"access_mode,omitempty"`
}

type VolumeManifestSpecCapacity struct {
	Storage string `json:"storage,omitempty" yaml:"storage,omitempty"`
}

func (v *VolumeManifest) FromJson(data []byte) error {
	return json.Unmarshal(data, v)
}

func (v *VolumeManifest) ToJson() ([]byte, error) {
	return json.Marshal(v)
}

func (v *VolumeManifest) FromYaml(data []byte) error {
	return yaml.Unmarshal(data, v)
}

func (v *VolumeManifest) ToYaml() ([]byte, error) {
	return yaml.Marshal(v)
}

func (v *VolumeManifest) SetVolumeMeta(vol *types.Volume) {

	if vol.Meta.Name == types.EmptyString {
		vol.Meta.Name = *v.Meta.Name
	}

	if v.Meta.Description != nil {
		vol.Meta.Description = *v.Meta.Description
	}

	if v.Meta.Labels != nil {
		vol.Meta.Labels = v.Meta.Labels
	}

}

func (v *VolumeManifest) SetVolumeSpec(vol *types.Volume) {

	t := vol.Spec.Updated
	defer func () {
		if t.Before(vol.Spec.Updated) {
			vol.Status.State = types.StateProvision
			return
		}
	}()

	if vol.Spec.Type != v.Spec.Type {
		vol.Spec.Type = v.Spec.Type
		vol.Spec.Updated = time.Now()
	}

	if vol.Spec.Selector.Node != v.Spec.Selector.Node {
		vol.Spec.Selector.Node = v.Spec.Selector.Node
		vol.Spec.Updated = time.Now()
	}

	var (
		ll = len(vol.Spec.Selector.Labels)
		lc = 0
	)

	for l, d := range vol.Spec.Selector.Labels {
		if _, ok := v.Spec.Selector.Labels[l]; !ok {
			continue
		}
		if v.Spec.Selector.Labels[l] != d {
			continue
		}
		lc++
	}

	if ll != lc {
		vol.Spec.Selector.Labels = make(map[string]string, 0)
		for l, d := range v.Spec.Selector.Labels {
			vol.Spec.Selector.Labels[l] = d
		}
		vol.Spec.Updated = time.Now()
	}

	stg, err := resource.DecodeMemoryResource(v.Spec.Capacity.Storage)
	if err != nil {
		return
	}

	if vol.Spec.Capacity.Storage != stg {
		vol.Spec.Capacity.Storage = stg
		vol.Spec.Updated = time.Now()
	}

	if vol.Spec.AccessMode != v.Spec.AccessMode {
		vol.Spec.AccessMode = v.Spec.AccessMode
		vol.Spec.Updated = time.Now()
	}

}


func (m VolumeManifest) GetManifest() *types.VolumeManifest {
	var v = new(types.VolumeManifest)

	v.Selector = m.Spec.Selector.GetSpec()
	v.Type = m.Spec.Type
	v.Capacity.Storage, _ = resource.DecodeMemoryResource(m.Spec.Capacity.Storage)
	v.AccessMode = m.Spec.AccessMode

	return v
}

type VolumeRemoveOptions struct {
	Force bool
}
