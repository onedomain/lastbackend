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
// +build linux

package cni

import (
	"github.com/onedomain/lastbackend/pkg/runtime/cni"
	"github.com/onedomain/lastbackend/pkg/runtime/cni/local"
	"github.com/onedomain/lastbackend/pkg/runtime/cni/vxlan"
	"github.com/spf13/viper"
)

func New(v *viper.Viper) (cni.CNI, error) {
	switch v.GetString("network.cni.type") {
	case "vxlan":
		return vxlan.New(v.GetString("network.interface"))
	default:
		return local.New()
	}
}
