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

package service

import (
	"context"
	"fmt"
	"github.com/onedomain/lastbackend/pkg/api/envs"
	"github.com/onedomain/lastbackend/pkg/api/types/v1/request"
	"github.com/onedomain/lastbackend/pkg/distribution"
	"github.com/onedomain/lastbackend/pkg/distribution/errors"
	"github.com/onedomain/lastbackend/pkg/distribution/types"
	"github.com/onedomain/lastbackend/pkg/log"
	"github.com/onedomain/lastbackend/pkg/util/resource"
	"net/http"
	"strings"
	"time"
)

const (
	logPrefix = "api:handler:service"
	logLevel  = 3
)

func Fetch(ctx context.Context, namespace, name string) (*types.Service, *errors.Err) {

	nm := distribution.NewServiceModel(ctx, envs.Get().GetStorage())
	svc, err := nm.Get(namespace, name)

	if err != nil {
		log.V(logLevel).Errorf("%s:fetch:> err: %s", logPrefix, err.Error())
		return nil, errors.New("service").InternalServerError(err)
	}

	if svc == nil {
		err := errors.New("service not found")
		log.V(logLevel).Errorf("%s:fetch:> err: %s", logPrefix, err.Error())
		return nil, errors.New("service").NotFound()
	}

	return svc, nil
}

func Apply(ctx context.Context, ns *types.Namespace, mf *request.ServiceManifest, opts *request.ServiceUpdateOptions) (*types.Service, *errors.Err) {

	if mf.Meta.Name == nil {
		return nil, errors.New("service").BadParameter("meta.name")
	}

	svc, err := Fetch(ctx, ns.Meta.Name, *mf.Meta.Name)
	if err != nil {
		if err.Code != http.StatusText(http.StatusNotFound) {
			return nil, errors.New("service").InternalServerError()
		}
	}

	if svc == nil {
		return Create(ctx, ns, mf)
	}

	return Update(ctx, ns, svc, mf, opts)
}

func Create(ctx context.Context, ns *types.Namespace, mf *request.ServiceManifest) (*types.Service, *errors.Err) {

	nm := distribution.NewNamespaceModel(ctx, envs.Get().GetStorage())
	sm := distribution.NewServiceModel(ctx, envs.Get().GetStorage())

	if mf.Meta.Name != nil {

		srv, err := sm.Get(ns.Meta.Name, *mf.Meta.Name)
		if err != nil {
			log.V(logLevel).Errorf("%s:create:> get service by name `%s` in namespace `%s` err: %s", logPrefix, mf.Meta.Name, ns.Meta.Name, err.Error())
			return nil, errors.New("service").InternalServerError()

		}

		if srv != nil {
			log.V(logLevel).Warnf("%s:create:> service name `%s` in namespace `%s` not unique", logPrefix, mf.Meta.Name, ns.Meta.Name)
			return nil, errors.New("service").NotUnique("name")

		}
	}

	svc := new(types.Service)
	mf.SetServiceMeta(svc)
	svc.Meta.SelfLink = *types.NewServiceSelfLink(ns.Meta.Name, *mf.Meta.Name)
	svc.Meta.Namespace = ns.Meta.Name
	svc.Meta.Endpoint = fmt.Sprintf("%s.%s", strings.ToLower(svc.Meta.Name), ns.Meta.Endpoint)

	if err := mf.SetServiceSpec(svc); err != nil {
		return nil, errors.New("service").BadRequest(err.Error())
	}

	if ns.Spec.Resources.Limits.RAM != 0 || ns.Spec.Resources.Limits.CPU != 0 {
		for _, c := range svc.Spec.Template.Containers {
			if c.Resources.Limits.RAM == 0 {
				c.Resources.Limits.RAM, _ = resource.DecodeMemoryResource(types.DEFAULT_RESOURCE_LIMITS_RAM)
			}
			if c.Resources.Limits.CPU == 0 {
				c.Resources.Limits.CPU, _ = resource.DecodeCpuResource(types.DEFAULT_RESOURCE_LIMITS_CPU)
			}
		}
	}

	if err := ns.AllocateResources(svc.Spec.GetResourceRequest()); err != nil {
		log.V(logLevel).Errorf("%s:create:> %s", logPrefix, err.Error())
		return nil, errors.New("service").BadRequest(err.Error())

	} else {
		if err := nm.Update(ns); err != nil {
			log.V(logLevel).Errorf("%s:update:> update namespace err: %s", logPrefix, err.Error())
			return nil, errors.New("service").InternalServerError()
		}
	}

	svc, err := sm.Create(ns, svc)
	if err != nil {
		log.V(logLevel).Errorf("%s:create:> create service err: %s", logPrefix, err.Error())
		return nil, errors.New("service").InternalServerError()
	}

	return svc, nil
}

func Update(ctx context.Context, ns *types.Namespace, svc *types.Service, mf *request.ServiceManifest, opts *request.ServiceUpdateOptions) (*types.Service, *errors.Err) {

	nm := distribution.NewNamespaceModel(ctx, envs.Get().GetStorage())
	sm := distribution.NewServiceModel(ctx, envs.Get().GetStorage())

	resources := svc.Spec.GetResourceRequest()

	mf.SetServiceMeta(svc)

	svc.Meta.Endpoint = fmt.Sprintf("%s.%s", strings.ToLower(svc.Meta.Name), ns.Meta.Endpoint)

	if err := mf.SetServiceSpec(svc); err != nil {
		return nil, errors.New("service").BadRequest(err.Error())
	}

	if opts.Redeploy {
		svc.Spec.Template.Updated = time.Now()
	}

	requestedResources := svc.Spec.GetResourceRequest()

	if !resources.Equal(requestedResources) {

		allocatedResources := ns.Status.Resources.Allocated
		ns.ReleaseResources(resources)

		if err := ns.AllocateResources(svc.Spec.GetResourceRequest()); err != nil {
			ns.Status.Resources.Allocated = allocatedResources
			log.V(logLevel).Errorf("%s:update:> %s", logPrefix, err.Error())
			return nil, errors.New("service").BadRequest(err.Error())
		} else {
			if err := nm.Update(ns); err != nil {
				log.V(logLevel).Errorf("%s:update:> update namespace err: %s", logPrefix, err.Error())
				return nil, errors.New("service").InternalServerError()
			}

		}
	}

	svc, err := sm.Update(svc)
	if err != nil {
		log.V(logLevel).Errorf("%s:update:> update service err: %s", logPrefix, err.Error())
		return nil, errors.New("service").InternalServerError()
	}

	return svc, nil
}
