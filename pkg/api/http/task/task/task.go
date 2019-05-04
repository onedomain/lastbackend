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

package task

import (
	"context"
	"github.com/onedomain/lastbackend/pkg/api/envs"
	"github.com/onedomain/lastbackend/pkg/api/types/v1/request"
	"github.com/onedomain/lastbackend/pkg/distribution"
	"github.com/onedomain/lastbackend/pkg/distribution/errors"
	"github.com/onedomain/lastbackend/pkg/distribution/types"
	"github.com/onedomain/lastbackend/pkg/log"
	"github.com/onedomain/lastbackend/pkg/util/generator"
	"github.com/onedomain/lastbackend/pkg/util/resource"
	"strings"
)

const (
	logPrefix = "api:handler:task"
	logLevel  = 3
)

func Fetch(ctx context.Context, namespace, job, name string) (*types.Task, *errors.Err) {

	tm := distribution.NewTaskModel(ctx, envs.Get().GetStorage())
	task, err := tm.Get(types.NewTaskSelfLink(namespace, job, name).String())

	if err != nil {
		log.V(logLevel).Errorf("%s:fetch:> err: %s", logPrefix, err.Error())
		return nil, errors.New("task").InternalServerError(err)
	}

	if task == nil {
		err := errors.New("task not found")
		log.V(logLevel).Errorf("%s:fetch:> err: %s", logPrefix, err.Error())
		return nil, errors.New("task").NotFound()
	}

	return task, nil
}

func Create(ctx context.Context, ns *types.Namespace, job *types.Job, mf *request.TaskManifest) (*types.Task, *errors.Err) {

	jm := distribution.NewJobModel(ctx, envs.Get().GetStorage())
	tm := distribution.NewTaskModel(ctx, envs.Get().GetStorage())

	if mf.Meta.Name != nil {

		task, err := tm.Get(types.NewTaskSelfLink(ns.Meta.Name, job.Meta.Name, *mf.Meta.Name).String())
		if err != nil {
			log.V(logLevel).Errorf("%s:create:> get task by name `%s` in namespace `%s` err: %s", logPrefix, mf.Meta.Name, ns.Meta.Name, err.Error())
			return nil, errors.New("task").InternalServerError()

		}

		if task != nil {
			log.V(logLevel).Warnf("%s:create:> task name `%s` in namespace `%s` not unique", logPrefix, mf.Meta.Name, ns.Meta.Name)
			return nil, errors.New("job").NotUnique("name")

		}
	}

	task := new(types.Task)
	task.Meta.SetDefault()
	task.Meta.Namespace = ns.Meta.Name
	task.Meta.Job = job.Meta.Name

	if mf.Meta.Name != nil {
		task.Meta.SelfLink = *types.NewTaskSelfLink(ns.Meta.Name, job.Meta.Name, *mf.Meta.Name)
		mf.SetTaskMeta(task)
	} else {
		name := strings.Split(generator.GetUUIDV4(), "-")[4][5:]
		task.Meta.Name = name
		task.Meta.SelfLink = *types.NewTaskSelfLink(ns.Meta.Name, job.Meta.Name, name)
	}

	task.Status.State = types.StateCreated

	task.Spec.Runtime = job.Spec.Task.Runtime
	task.Spec.Selector = job.Spec.Task.Selector
	task.Spec.Template = job.Spec.Task.Template

	if err := mf.SetTaskSpec(task); err != nil {
		log.V(logLevel).Errorf("%s:create:> set task spec err: %s", logPrefix, err.Error())
		return nil, errors.New("task").BadParameter("spec")
	}

	if job.Spec.Resources.Limits.RAM != 0 || job.Spec.Resources.Limits.CPU != 0 {
		for _, c := range task.Spec.Template.Containers {
			if c.Resources.Limits.RAM == 0 {
				c.Resources.Limits.RAM, _ = resource.DecodeMemoryResource(types.DEFAULT_RESOURCE_LIMITS_RAM)
			}
			if c.Resources.Limits.CPU == 0 {
				c.Resources.Limits.CPU, _ = resource.DecodeCpuResource(types.DEFAULT_RESOURCE_LIMITS_CPU)
			}
		}
	}

	if err := job.AllocateResources(task.Spec.GetResourceRequest()); err != nil {
		log.V(logLevel).Errorf("%s:create:> %s", logPrefix, err.Error())
		return nil, errors.New("job").BadRequest(err.Error())
	} else {
		if err := jm.Set(job); err != nil {
			log.V(logLevel).Errorf("%s:update:> update namespace err: %s", logPrefix, err.Error())
			return nil, errors.New("job").InternalServerError()
		}
	}

	if _, err := tm.Create(task); err != nil {
		log.V(logLevel).Errorf("%s:create:> create task err: %s", logPrefix, err.Error())
		return nil, errors.New("task").InternalServerError()
	}

	return task, nil
}
