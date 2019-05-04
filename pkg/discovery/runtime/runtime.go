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
	"context"
)

const (
	logRuntimePrefix = "discovery:runtime"
	logLevel         = 3
)

type Runtime struct {
	ctx context.Context

	opts *RuntimeOpts
}

type RuntimeOpts struct {
	Iface string
	Port  uint16
}

func New(opts *RuntimeOpts) *Runtime {
	return &Runtime{
		ctx:  context.Background(),
		opts: opts,
	}
}
