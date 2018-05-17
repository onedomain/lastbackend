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

package endpoint

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
)

const logEndpointPrefix = "runtime:endpoint:>"

func Restore(ctx context.Context) error {
	log.Debugf("%s restore", logEndpointPrefix)
	return nil
}

func Create(ctx context.Context, key string, spec *types.EndpointSpec) (types.EndpointStatus, error) {
	log.Debugf("%s create", logEndpointPrefix)
	return types.EndpointStatus{}, nil
}

func Clean(ctx context.Context, status *types.EndpointStatus) {
	log.Debugf("%s clean", logEndpointPrefix)
}

func Destroy(ctx context.Context, endpoint string, status *types.EndpointStatus) {
	log.Debugf("%s destroy", logEndpointPrefix)
}

func Manage(ctx context.Context) error {
	log.Debugf("%s manage", logEndpointPrefix)
	return nil
}