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

package cpi

import (
	"context"
	"github.com/onedomain/lastbackend/pkg/distribution/types"
)

type CPI interface {
	Info(ctx context.Context) (map[string]*types.EndpointState, error)
	Create(ctx context.Context, manifest *types.EndpointManifest) (*types.EndpointState, error)
	Destroy(ctx context.Context, state *types.EndpointState) error
	Update(ctx context.Context, state *types.EndpointState, manifest *types.EndpointManifest) (*types.EndpointState, error)
}
