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

package network

import (
	"github.com/onedomain/lastbackend/pkg/distribution/types"
	"github.com/onedomain/lastbackend/pkg/log"
	"github.com/onedomain/lastbackend/pkg/network/state"
	"github.com/onedomain/lastbackend/pkg/runtime/cni"
	ni "github.com/onedomain/lastbackend/pkg/runtime/cni/cni"
	"github.com/onedomain/lastbackend/pkg/runtime/cpi"
	pi "github.com/onedomain/lastbackend/pkg/runtime/cpi/cpi"
	"github.com/spf13/viper"
)

const (
	DefaultResolverIP = "172.17.0.1"
)

type Network struct {
	state    *state.State
	cni      cni.CNI
	cpi      cpi.CPI
	resolver struct {
		ip       string
		external []string
	}
}

func New(v *viper.Viper) (*Network, error) {

	var err error

	net := new(Network)

	if v.GetString("runtime.cni.type") == types.EmptyString &&
		v.GetString("runtime.cpi.type") == types.EmptyString {
		log.Debug("run without network management")
		return nil, nil
	}

	net.state = state.New()
	if net.cni, err = ni.New(v); err != nil {
		log.Errorf("Cannot initialize cni: %v", err)
		return nil, err
	}

	if net.cpi, err = pi.New(v); err != nil {
		log.Errorf("Cannot initialize cni: %v", err)
		return nil, err
	}

	rip := v.GetString("network.resolver.ip")
	if rip == types.EmptyString {
		rip = DefaultResolverIP
	}

	net.resolver.ip = rip
	net.resolver.external = v.GetStringSlice("network.resolver.servers")
	if len(net.resolver.external) == 0 {
		net.resolver.external = []string{"8.8.8.8", "8.8.4.4"}
	}

	return net, nil
}
