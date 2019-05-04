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

import "fmt"

type Image struct {
	Meta   ImageMeta
	Status ImageStatus
	Spec   ImageSpec
}

type ImageMeta struct {
	ID   string   `json:"id"`
	Digest string `json:"digest"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type ImageStatus struct {
	State       string `json:"state"`
	Size        int64  `json:"size"`
	VirtualSize int64  `json:"virtual_size"`
	Container   ImageContainer `json:"container"`
}

type ImageContainer struct {
	Ports []string
	Envs  []string
	Exec  SpecTemplateContainerExec
}

type ImageSpec struct {
	// Name full name
	Name string `json:"name"`
	// Secret name for pulling
	Secret string `json:"auth"`
}

type ImageManifest struct {
	Name   string `json:"name" yaml:"name"`
	Tag    string `json:"tag" yaml:"tag"`
	Auth   string `json:"auth" yaml:"auth"`
	Policy string `json:"policy" yaml:"policy"`
}

func (i *Image) SelfLink() string {
	return fmt.Sprintf("%s", i.Meta.Name)
}
