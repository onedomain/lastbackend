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

package job

import (
	"context"
	"github.com/onedomain/lastbackend/pkg/util/generator"
	"strings"

	"github.com/onedomain/lastbackend/pkg/controller/envs"
	"github.com/onedomain/lastbackend/pkg/distribution"
	"github.com/onedomain/lastbackend/pkg/distribution/errors"
	"github.com/onedomain/lastbackend/pkg/distribution/types"
	"github.com/onedomain/lastbackend/pkg/log"
	"time"
)

const logPodPrefix = "state:observer:pod"

// PodObserve function manages pod handlers based on pod state
func PodObserve(js *JobState, p *types.Pod) (err error) {

	log.V(logLevel).Debugf("%s:> observe start: %s > state %s", logPodPrefix, p.SelfLink(), p.Status.State)

	// Call pod state manager methods
	switch p.Status.State {
	case types.StateCreated:
		err = handlePodStateCreated(js, p)
	case types.StateProvision:
		err = handlePodStateProvision(js, p)
	case types.StateReady:
		err = handlePodStateReady(js, p)
	case types.StateError:
		err = handlePodStateError(js, p)
	case types.StateDegradation:
		err = handlePodStateDegradation(js, p)
	case types.StateDestroy:
		err = handlePodStateDestroy(js, p)
	case types.StateDestroyed:
		err = handlePodStateDestroyed(js, p)
	}
	if err != nil {
		log.Errorf("%s:> handle pod state %s err: %s", logPodPrefix, p.Status.State, err.Error())
		return err
	}

	log.V(logLevel).Debugf("%s:> observe state finish: %s", logPodPrefix, p.SelfLink())

	_, sl := p.SelfLink().Parent()
	if p.Status.State != types.StateDestroyed {
		js.pod.list[sl.String()] = p
	}

	d, ok := js.task.list[sl.String()]
	if !ok {
		log.V(logLevel).Debugf("%s:> task node found: %s", logPodPrefix, sl.String())
		return nil
	}

	log.V(logLevel).Debugf("%s:> observe finish: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)
	if err := taskStatusState(js, d, p); err != nil {
		return err
	}

	return nil
}

func handlePodStateCreated(ss *JobState, p *types.Pod) error {

	log.V(logLevel).Debugf("%s:> handlePodStateCreated: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	if err := podProvision(ss, p); err != nil {
		return err
	}
	log.V(logLevel).Debugf("%s handle pod create state finish: %s : %s", logPodPrefix, p.SelfLink(), p.Status.State)
	return nil
}

func handlePodStateProvision(ss *JobState, p *types.Pod) error {

	log.V(logLevel).Debugf("%s:> handlePodStateProvision: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	if err := podProvision(ss, p); err != nil {
		return err
	}
	return nil
}

func handlePodStateReady(ss *JobState, p *types.Pod) error {

	log.V(logLevel).Debugf("%s:> handlePodStateReady: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	return nil
}

func handlePodStateError(ss *JobState, p *types.Pod) error {

	log.V(logLevel).Debugf("%s:> handlePodStateError: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	return nil
}

func handlePodStateDegradation(ss *JobState, p *types.Pod) error {

	log.V(logLevel).Debugf("%s:> handlePodStateDegradation: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	return nil
}

func handlePodStateDestroy(ss *JobState, p *types.Pod) error {

	log.V(logLevel).Debugf("%s:> handlePodStateDestroy: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	if err := podDestroy(ss, p); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

func handlePodStateDestroyed(ss *JobState, p *types.Pod) error {

	log.V(logLevel).Debugf("%s:> handlePodStateDestroyed: %s > %s", logPodPrefix, p.SelfLink(), p.Status.State)

	if err := podRemove(ss, p); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	return nil
}

// podCreate function creates new pod based on task spec
func podCreate(d *types.Task) (*types.Pod, error) {
	dm := distribution.NewPodModel(context.Background(), envs.Get().GetStorage())

	pod := types.NewPod()
	pod.Meta.SetDefault()
	pod.Meta.Name = strings.Split(generator.GetUUIDV4(), "-")[4][5:]
	pod.Meta.Namespace = d.Meta.Namespace
	sl, _ := types.NewPodSelfLink(types.KindTask, d.SelfLink().String(), pod.Meta.Name)
	pod.Meta.SelfLink = *sl
	pod.Status.SetCreated()

	pod.Spec.SetSpecRuntime(d.Spec.Runtime)
	pod.Spec.SetSpecTemplate(pod.SelfLink().String(), d.Spec.Template)
	pod.Spec.Selector = d.Spec.Selector

	return dm.Put(pod)
}

// podDestroy function marks pod as provision
func podProvision(ss *JobState, p *types.Pod) (err error) {

	t := p.Meta.Updated

	defer func() {

		if err == nil {
			err = podUpdate(p, t)
		}

	}()

	if p.Meta.Node == types.EmptyString {

		var node *types.Node

		node, err = ss.cluster.PodLease(p)
		if err != nil {
			log.Errorf("%s:> pod node lease err: %s", logPrefix, err.Error())
			return err
		}

		if node == nil {
			p.Status.State = types.StateError
			p.Status.Message = errors.NodeNotFound
			p.Meta.Updated = time.Now()
			return nil
		}

		p.Meta.Node = node.SelfLink().String()
		p.Meta.Updated = time.Now()
	}

	if err = podManifestPut(p); err != nil {
		log.Errorf("%s:> pod manifest create err: %s", logPrefix, err.Error())
		return err
	}

	if p.Status.State != types.StateProvision {
		p.Status.State = types.StateProvision
		p.Meta.Updated = time.Now()
	}

	return nil
}

// podDestroy function marks pod spec as destroy
func podDestroy(ss *JobState, p *types.Pod) (err error) {

	t := p.Meta.Updated
	defer func() {
		if err == nil {
			err = podUpdate(p, t)
		}
	}()

	if p.Spec.State.Destroy {

		if p.Meta.Node == types.EmptyString {
			p.Status.State = types.StateDestroyed
			p.Meta.Updated = time.Now()
			return nil
		}

		if p.Status.State != types.StateDestroy {
			p.Status.State = types.StateDestroy
			p.Meta.Updated = time.Now()
		}

		return nil
	}

	p.Spec.State.Destroy = true
	if err = podManifestSet(p); err != nil {
		if errors.Storage().IsErrEntityNotFound(err) {
			if p.Meta.Node != types.EmptyString {
				if _, err := ss.cluster.PodRelease(p); err != nil {
					if !errors.Storage().IsErrEntityNotFound(err) {
						return err
					}
				}
			}

			p.Status.State = types.StateDestroyed
			p.Meta.Updated = time.Now()
			return nil
		}

		return err
	}

	p.Status.State = types.StateDestroy

	if p.Meta.Node == types.EmptyString {
		p.Status.State = types.StateDestroyed
	}

	p.Meta.Updated = time.Now()
	return nil
}

// podRemove function removes pod from storage if node is released
func podRemove(ss *JobState, p *types.Pod) (err error) {

	pm := distribution.NewPodModel(context.Background(), envs.Get().GetStorage())
	if _, err = ss.cluster.PodRelease(p); err != nil {
		return err
	}

	if err = podManifestDel(p); err != nil {
		return err
	}

	p.Meta.Node = types.EmptyString
	p.Meta.Updated = time.Now()

	if err = pm.Remove(p); err != nil {
		log.Errorf("%s", err.Error())
		return err
	}

	ss.DelPod(p)
	return nil
}

func podUpdate(p *types.Pod, timestamp time.Time) error {

	if timestamp.Before(p.Meta.Updated) {
		pm := distribution.NewPodModel(context.Background(), envs.Get().GetStorage())
		if err := pm.Update(p); err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
	}

	return nil
}

func podManifestPut(p *types.Pod) error {

	mm := distribution.NewPodModel(context.Background(), envs.Get().GetStorage())
	m, err := mm.ManifestGet(p.Meta.Node, p.Meta.SelfLink.String())
	if err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			log.Errorf("%s", err.Error())
			return err
		}
	}

	if m == nil {
		pm := types.PodManifest(p.Spec)

		if err := mm.ManifestAdd(p.Meta.Node, p.Meta.SelfLink.String(), &pm); err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
	}

	return nil
}

func podManifestSet(p *types.Pod) error {

	var (
		m   *types.PodManifest
		err error
	)

	mm := distribution.NewPodModel(context.Background(), envs.Get().GetStorage())
	m, err = mm.ManifestGet(p.Meta.Node, p.Meta.SelfLink.String())
	if err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			log.Errorf("%s", err.Error())
			return err
		}
	}

	// Update manifest
	if m == nil {
		ms := types.PodManifest(p.Spec)
		m = &ms
	} else {
		*m = types.PodManifest(p.Spec)
	}

	if err := mm.ManifestSet(p.Meta.Node, p.Meta.SelfLink.String(), m); err != nil {
		return err
	}

	return nil
}

func podManifestDel(p *types.Pod) error {

	if p.Meta.Node == types.EmptyString {
		return nil
	}

	// Remove manifest
	mm := distribution.NewPodModel(context.Background(), envs.Get().GetStorage())
	err := mm.ManifestDel(p.Meta.Node, p.SelfLink().String())
	if err != nil {
		if !errors.Storage().IsErrEntityNotFound(err) {
			return err
		}
	}

	return nil
}
