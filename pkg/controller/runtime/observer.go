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
	"github.com/onedomain/lastbackend/pkg/controller/envs"
	"github.com/onedomain/lastbackend/pkg/controller/state"
	"github.com/onedomain/lastbackend/pkg/storage"
	"golang.org/x/net/context"
)

const (
	logPrefix = "controller:runtime:observer"
)

type Observer struct {
	rev   *int64
	stg   storage.Storage
	state *state.State
}

func (o *Observer) Loop() {
	o.state.Loop()
}

func (o *Observer) Stop() {
	o.state = nil
}

func NewObserver(ctx context.Context) *Observer {

	o := new(Observer)

	o.stg = envs.Get().GetStorage()
	o.state = state.NewState()

	return o
}
