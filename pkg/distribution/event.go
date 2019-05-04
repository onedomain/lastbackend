//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package distribution

import (
	"context"
	"github.com/onedomain/lastbackend/pkg/distribution/types"
	"github.com/onedomain/lastbackend/pkg/log"
	"github.com/onedomain/lastbackend/pkg/storage"
	//"regexp"
)

const (
	logEventPrefix = "distribution:events"
)

type Event struct {
	context context.Context
	storage storage.Storage
}

func (e *Event) Runtime() (*types.System, error) {

	log.V(logLevel).Debugf("%s:get:> get events runtime info", logEventPrefix)
	runtime, err := e.storage.Info(e.context, e.storage.Collection().Root(), "")
	if err != nil {
		log.V(logLevel).Errorf("%s:get:> get runtime info error: %s", logEventPrefix, err)
		return &runtime.System, err
	}
	return &runtime.System, nil
}
