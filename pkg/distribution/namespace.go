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

package distribution

import (
	"context"
	"encoding/json"
	"github.com/onedomain/lastbackend/pkg/distribution/errors"
	"github.com/onedomain/lastbackend/pkg/distribution/types"
	"github.com/onedomain/lastbackend/pkg/log"
	"github.com/onedomain/lastbackend/pkg/storage"
)

const (
	logNamespacePrefix  = "distribution:namespace"
	defaultNamespaceRam = "2GB"
	defaultNamespaceCPU = "200m"
)

type Namespace struct {
	context context.Context
	storage storage.Storage
}

type NM struct {
	Meta   struct{}
	Entity Namespace
}

func (n *NM) Set(Namespace) error {
	return nil
}

func (n *Namespace) List() (*types.NamespaceList, error) {

	log.V(logLevel).Debugf("%s:list:> get namespaces list", logNamespacePrefix)

	var list = types.NewNamespaceList()

	err := n.storage.List(n.context, n.storage.Collection().Namespace(), "", list, nil)

	if err != nil {
		log.Info(err.Error())
		log.V(logLevel).Error("%s:list:> get namespaces list err: %v", logNamespacePrefix, err)
		return nil, err
	}

	log.V(logLevel).Debugf("%s:list:> get namespaces list result: %d", logNamespacePrefix, len(list.Items))

	return list, nil
}

func (n *Namespace) Get(name string) (*types.Namespace, error) {

	log.V(logLevel).Infof("%s:get:> get namespace %s", logNamespacePrefix, name)

	if name == "" {
		return nil, errors.New(errors.ArgumentIsEmpty)
	}

	namespace := new(types.Namespace)
	key := types.NewNamespaceSelfLink(name).String()

	err := n.storage.Get(n.context, n.storage.Collection().Namespace(), key, &namespace, nil)
	if err != nil {
		if errors.Storage().IsErrEntityNotFound(err) {
			log.V(logLevel).Warnf("%s:get:> namespace by name `%s` not found", logNamespacePrefix, name)
			return nil, nil
		}

		log.V(logLevel).Errorf("%s:get:> get namespace by name `%s` err: %v", logNamespacePrefix, name, err)
		return nil, err
	}

	return namespace, nil
}

func (n *Namespace) Create(ns *types.Namespace) (*types.Namespace, error) {

	log.V(logLevel).Debugf("%s:create:> create Namespace %#v", logNamespacePrefix, ns.Meta.Name)

	ns.Meta.SetDefault()
	ns.Meta.SelfLink = *types.NewNamespaceSelfLink(ns.Meta.Name)

	if err := n.storage.Put(n.context, n.storage.Collection().Namespace(), ns.SelfLink().String(), ns, nil); err != nil {
		log.V(logLevel).Errorf("%s:create:> insert namespace err: %v", logNamespacePrefix, err)
		return nil, err
	}

	return ns, nil
}

func (n *Namespace) Update(namespace *types.Namespace) error {

	log.V(logLevel).Debugf("%s:update:> update Namespace %#v", logNamespacePrefix, namespace)

	if err := n.storage.Set(n.context, n.storage.Collection().Namespace(),
		namespace.SelfLink().String(), namespace, nil); err != nil {
		log.V(logLevel).Errorf("%s:update:> namespace update err: %v", logNamespacePrefix, err)
		return err
	}

	return nil
}

func (n *Namespace) Remove(ns *types.Namespace) error {

	log.V(logLevel).Debugf("%s:remove:> remove namespace %s", logNamespacePrefix, ns.Meta.Name)

	if err := n.storage.Del(n.context, n.storage.Collection().Namespace(), ns.SelfLink().String()); err != nil {
		log.V(logLevel).Errorf("%s:remove:> remove namespace err: %v", logNamespacePrefix, err)
		return err
	}

	return nil
}

// Watch namespace changes
func (n *Namespace) Watch(ch chan types.NamespaceEvent) error {

	log.V(logLevel).Debugf("%s:watch:> watch namespace", logNamespacePrefix)

	done := make(chan bool)
	watcher := storage.NewWatcher()

	go func() {
		for {
			select {
			case <-n.context.Done():
				done <- true
				return
			case e := <-watcher:
				if e.Data == nil {
					continue
				}

				res := types.NamespaceEvent{}
				res.Action = e.Action
				res.Name = e.Name

				obj := new(types.Namespace)

				if err := json.Unmarshal(e.Data.([]byte), &obj); err != nil {
					log.Errorf("%s:watch:> parse json", logNamespacePrefix)
					continue
				}

				res.Data = obj

				ch <- res
			}
		}
	}()

	opts := storage.GetOpts()
	if err := n.storage.Watch(n.context, n.storage.Collection().Namespace(), watcher, opts); err != nil {
		return err
	}

	return nil
}

func NewNamespaceModel(ctx context.Context, stg storage.Storage) *Namespace {
	return &Namespace{ctx, stg}
}
