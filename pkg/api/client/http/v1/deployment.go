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

	"github.com/onedomain/lastbackend/pkg/api/client/types"
	rv1 "github.com/onedomain/lastbackend/pkg/api/types/v1/request"
	vv1 "github.com/onedomain/lastbackend/pkg/api/types/v1/views"
	"github.com/onedomain/lastbackend/pkg/distribution/errors"
	t "github.com/onedomain/lastbackend/pkg/distribution/types"
	"github.com/onedomain/lastbackend/pkg/util/http/request"
)

type DeploymentClient struct {
	client *request.RESTClient

	namespace t.NamespaceSelfLink
	service   t.ServiceSelfLink
	selflink  t.DeploymentSelfLink
}

func (dc *DeploymentClient) Pod(args ...string) types.PodClientV1 {
	name := ""
	// Get any parameters passed to us out of the args variable into "real"
	// variables we created for them.
	for i := range args {
		switch i {
		case 0: // hostname
			name = args[0]
		default:
			panic("Wrong parameter count: (is allowed from 0 to 1)")
		}
	}
	return newPodClient(dc.client, dc.namespace.String(), t.KindDeployment, dc.selflink.String(), name)
}

func (dc *DeploymentClient) List(ctx context.Context) (*vv1.DeploymentList, error) {

	var s *vv1.DeploymentList
	var e *errors.Http

	err := dc.client.Get(fmt.Sprintf("/namespace/%s/service/%s/deploymet", dc.namespace.String(), dc.service.String())).
		AddHeader("Content-Type", "application/json").
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	if s == nil {
		list := make(vv1.DeploymentList, 0)
		s = &list
	}

	return s, nil
}

func (dc *DeploymentClient) Get(ctx context.Context) (*vv1.Deployment, error) {

	var s *vv1.Deployment
	var e *errors.Http

	err := dc.client.Get(fmt.Sprintf("/namespace/%s/service/%s/deploymet/%s", dc.namespace.String(), dc.service.String(), dc.selflink.Name())).
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

func (dc *DeploymentClient) Update(ctx context.Context, opts *rv1.DeploymentUpdateOptions) (*vv1.Deployment, error) {

	body, err := opts.ToJson()
	if err != nil {
		return nil, err
	}

	var s *vv1.Deployment
	var e *errors.Http

	err = dc.client.Put(fmt.Sprintf("/namespace/%s/service/%s/deployment/%s", dc.namespace.String(), dc.service.String(), dc.selflink.Name())).
		AddHeader("Content-Type", "application/json").
		Body(body).
		JSON(&s, &e)

	if err != nil {
		return nil, err
	}
	if e != nil {
		return nil, errors.New(e.Message)
	}

	return s, nil
}

func newDeploymentClient(client *request.RESTClient, namespace, service, name string) *DeploymentClient {
	return &DeploymentClient{
		client:    client,
		namespace: *t.NewNamespaceSelfLink(namespace),
		service:   *t.NewServiceSelfLink(namespace, service),
		selflink:  *t.NewDeploymentSelfLink(namespace, service, name),
	}
}
