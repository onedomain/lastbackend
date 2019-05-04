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

import (
	"time"
)

const (
	KindVolumeHostDir = "dir"
	KindVolumeLocal   = "local"
)

// swagger:ignore
// swagger:model types_volume
type Volume struct {
	System
	// Volume meta
	Meta VolumeMeta `json:"meta" yaml:"meta"`
	// Volume spec
	Spec VolumeSpec `json:"spec" yaml:"spec"`
	// Volume status
	Status VolumeStatus `json:"status" yaml:"status"`
}

// swagger:ignore
// swagger:model types_volume_map
type VolumeMap struct {
	System
	Items map[string]*Volume
}

// swagger:ignore
// swagger:model types_volume_list
type VolumeList struct {
	System
	Items []*Volume
}

// swagger:ignore
// swagger:model types_volume_meta
type VolumeMeta struct {
	Meta
	SelfLink  VolumeSelfLink `json:"self_link"`
	Node      string         `json:"node"`
	Namespace string         `json:"namespace"`
}

// swagger:model types_volume_spec
type VolumeSpec struct {
	Type       string             `json:"type"`
	Selector   SpecSelector       `json:"selector"`
	Capacity   SpecVolumeCapacity `json:"capacity"`
	State      VolumeSpecState    `json:"state"`
	HostPath   string             `json:"host_path"`
	AccessMode string             `json:"access_mode"`

	Updated time.Time `json:"updated"`
}

// swagger:model types_volume_spec_state
type VolumeSpecState struct {
	Destroy bool `json:"destroy"`
}

// swagger:ignore
// swagger:model types_volume_create
type VolumeCreateOptions struct {
}

// swagger:ignore
// swagger:model types_volume_update
type VolumeUpdateOptions struct {
}

type VolumeStatus struct {
	// volume state
	State string `json:"state" yaml:"state"`
	// volume status
	Status VolumeState `json:"status" yaml:"status"`
	// volume status message
	Message string `json:"message" yaml:"message"`
}

// swagger:ignore
// swagger:model types_volume_status
type VolumeState struct {
	Type string `json:"type" yaml:"type"`
	// Volume root path
	Path string `json:"path" yaml:"path"`
	// Volume state ready
	Ready bool `json:"ready" yaml:"ready"`
}

func (vs *VolumeStatus) SetReady() {
	vs.Status.Ready = true
	vs.State = StateReady
	vs.Message = EmptyString
}

func (vs *VolumeStatus) SetDestroyed() {
	vs.Status.Ready = false
	vs.State = StateDestroyed
	vs.Message = EmptyString
}

func (vs *VolumeStatus) SetError(err error) {
	vs.Status.Ready = false
	vs.State = StateError
	vs.Message = err.Error()
}

func (v *Volume) SelfLink() *VolumeSelfLink {
	return &v.Meta.SelfLink
}

func NewVolumeList() *VolumeList {
	dm := new(VolumeList)
	dm.Items = make([]*Volume, 0)
	return dm
}

func NewVolumeMap() *VolumeMap {
	dm := new(VolumeMap)
	dm.Items = make(map[string]*Volume)
	return dm
}

func NewVolumeStatus() *VolumeStatus {
	status := VolumeStatus{}
	return &status
}
