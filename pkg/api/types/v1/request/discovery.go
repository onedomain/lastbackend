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

import "github.com/onedomain/lastbackend/pkg/distribution/types"

// swagger:model request_ingress_connect
type DiscoveryConnectOptions struct {
	Info    types.DiscoveryInfo   `json:"info"`
	Status  types.DiscoveryStatus `json:"status"`
	Network types.NetworkState    `json:"network"`
}

type DiscoveryMeta struct {
}

// swagger:ignore
// swagger:model request_node_remove
type DiscoveryRemoveOptions struct {
	Force bool `json:"force"`
}

type DiscoveryStatusOptions types.DiscoveryStatus
