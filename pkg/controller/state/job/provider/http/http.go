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

package http

import (
	"github.com/onedomain/lastbackend/pkg/api/types/v1/request"
	"github.com/onedomain/lastbackend/pkg/distribution/types"
	"github.com/onedomain/lastbackend/pkg/log"
	"net/http"
	"strings"
	"time"
)

const (
	logLevel = 3
)

type JobHttpProvider struct {
	timeout time.Time
	config  *types.JobSpecProviderHTTP
	client  http.Client
}

func (h *JobHttpProvider) Fetch() (*types.TaskManifest, error) {

	var (
		err      error
		manifest = new(request.TaskManifest)
	)

	client := http.Client{}

	req, err := http.NewRequest(strings.ToUpper(h.config.Method), h.config.Endpoint, nil)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	if len(h.config.Headers) > 0 {
		for k, v := range h.config.Headers {
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	if err := manifest.DecodeAndValidate(resp.Body); err != nil {
		log.Error(err.Err().Error())
		return nil, err.Err()
	}

	defer resp.Body.Close()

	mf := new(types.TaskManifest)
	manifest.SetTaskManifestMeta(mf)
	if err := manifest.SetTaskManifestSpec(mf); err != nil {
		return nil, err
	}

	return mf, nil
}

func New(cfg *types.JobSpecProviderHTTP) (*JobHttpProvider, error) {

	log.V(logLevel).Debug("Use http task watcher")

	var (
		provider *JobHttpProvider
	)

	provider = new(JobHttpProvider)
	provider.config = cfg

	return provider, nil
}
