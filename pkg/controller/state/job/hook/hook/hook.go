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

package hook

import (
	"github.com/onedomain/lastbackend/pkg/controller/state/job/hook"
	"github.com/onedomain/lastbackend/pkg/controller/state/job/hook/http"
	"github.com/onedomain/lastbackend/pkg/distribution/types"
)

func New(hook types.JobSpecHook) (hook.Hook, error) {

	if hook.Http != nil {
		return http.New(hook.Http)
	}

	return nil, nil

}
