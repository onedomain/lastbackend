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

package v1

import (
	"context"
	"fmt"
	rv1 "github.com/onedomain/lastbackend/pkg/api/types/v1/request"
	vv1 "github.com/onedomain/lastbackend/pkg/api/types/v1/views"
	"github.com/onedomain/lastbackend/pkg/distribution/errors"
	"github.com/onedomain/lastbackend/pkg/util/http/request"
)

type IngressClient struct {
	client *request.RESTClient

	hostname string
}

func (ic *IngressClient) List(ctx context.Context) (*vv1.IngressList, error) {

	var i *vv1.IngressList
	var e *errors.Http

	err := ic.client.Get(fmt.Sprintf("/ingress")).
		AddHeader("Content-Type", "application/json").
		JSON(&i, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	if ic == nil {
		list := make(vv1.IngressList, 0)
		i = &list
	}

	return i, nil
}

func (ic *IngressClient) Get(ctx context.Context) (*vv1.Ingress, error) {

	var s *vv1.Ingress
	var e *errors.Http

	err := ic.client.Get(fmt.Sprintf("/ingress/%s", ic.hostname)).
		AddHeader("Content-Type", "application/json").
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	return s, nil
}

func (ic *IngressClient) Connect(ctx context.Context, opts *rv1.IngressConnectOptions) error {

	body := opts.ToJson()

	var e *errors.Http

	err := ic.client.Put(fmt.Sprintf("/ingress/%s", ic.hostname)).
		AddHeader("Content-Type", "application/json").
		Body([]byte(body)).
		JSON(nil, &e)

	if err != nil {
		return err
	}
	if e != nil {
		return errors.New(e.Message)
	}

	return nil
}

func (ic *IngressClient) SetStatus(ctx context.Context, opts *rv1.IngressStatusOptions) (*vv1.IngressManifest, error) {

	body := opts.ToJson()

	var s *vv1.IngressManifest
	var e *errors.Http

	err := ic.client.Put(fmt.Sprintf("/ingress/%s/status", ic.hostname)).
		AddHeader("Content-Type", "application/json").
		Body([]byte(body)).
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	return s, nil
}

func newIngressClient(req *request.RESTClient, hostname string) *IngressClient {
	return &IngressClient{client: req, hostname: hostname}
}
