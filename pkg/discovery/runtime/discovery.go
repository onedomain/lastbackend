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

package runtime

import (
	"fmt"
	"os"

	"github.com/onedomain/lastbackend/pkg/distribution/types"
	"github.com/onedomain/lastbackend/pkg/util/system"
)

func (r *Runtime) DiscoveryInfo() types.DiscoveryInfo {

	var (
		info = types.DiscoveryInfo{}
	)

	osInfo := system.GetOsInfo()
	hostname, err := os.Hostname()
	if err != nil {
		_ = fmt.Errorf("get hostname err: %s", err)
	}

	ip, err := system.GetHostIP(r.opts.Iface)
	if err != nil {
		_ = fmt.Errorf("get ip err: %s", err)
	}

	info.Hostname = hostname
	info.InternalIP = ip
	info.OSType = osInfo.GoOS
	info.OSName = fmt.Sprintf("%s %s", osInfo.OS, osInfo.Core)
	info.Architecture = osInfo.Platform

	return info
}

func (r *Runtime) DiscoveryStatus() types.DiscoveryStatus {

	var state = types.DiscoveryStatus{}

	ip, err := system.GetHostIP(r.opts.Iface)
	if err != nil {
		_ = fmt.Errorf("get ip err: %s", err)
	}

	state.Port = r.opts.Port
	state.IP = ip

	return state
}
