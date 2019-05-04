//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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

package etcd

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"strings"

	"github.com/onedomain/lastbackend/pkg/log"
	"github.com/onedomain/lastbackend/pkg/storage/etcd/store"
	"github.com/onedomain/lastbackend/pkg/storage/etcd/v3"
	"github.com/onedomain/lastbackend/pkg/storage/types"
)

const (
	logLevel  = 6
	logPrefix = "storage:etcd"
)

type Storage struct {
	client *client
}

type client struct {
	store store.Store
	dfunc store.DestroyFunc
}

func New(config *v3.Config) (*Storage, error) {

	log.V(logLevel).Debug("Etcd: define storage")

	var (
		err    error
		s      = new(Storage)
	)

	s.client = new(client)

	if s.client.store, s.client.dfunc, err = v3.GetClient(config); err != nil {
		log.Errorf("%s: store initialize err: %v", logPrefix, err)
		return nil, err
	}

	return s, nil
}

func (s Storage) Info(ctx context.Context, collection string, name string) (*types.System, error) {
	return s.client.store.Info(ctx, keyCreate(collection, name))
}

func (s Storage) Get(ctx context.Context, collection string, name string, obj interface{}, opts *types.Opts) error {
	if reflect.ValueOf(obj).IsNil() {
		return errors.New(types.ErrStructOutIsNil)
	}

	var rev *int64
	if opts != nil {
		rev = opts.Rev
	}

	return s.client.store.Get(ctx, keyCreate(collection, name), obj, rev)
}

func (s Storage) List(ctx context.Context, collection string, query string, obj interface{}, opts *types.Opts) error {

	if reflect.ValueOf(obj).IsNil() {
		return errors.New(types.ErrStructOutIsNil)
	}

	var rev *int64
	if opts != nil {
		rev = opts.Rev
	}

	return s.client.store.List(ctx, keyCreate(collection, query), "", obj, rev)
}

func (s Storage) Map(ctx context.Context, collection string, query string, obj interface{}, opts *types.Opts) error {

	if reflect.ValueOf(obj).IsNil() {
		return errors.New(types.ErrStructOutIsNil)
	}

	q := ".*/(.*)$"

	var rev *int64
	if opts != nil {
		rev = opts.Rev
	}

	return s.client.store.Map(ctx, keyCreate(collection, query), q, obj, rev)
}

func (s Storage) Put(ctx context.Context, collection string, name string, obj interface{}, opts *types.Opts) error {

	if opts == nil {
		opts = new(types.Opts)
	}

	return s.client.store.Put(ctx, keyCreate(collection, name), obj, nil, opts.Ttl)
}

func (s Storage) Set(ctx context.Context, collection string, name string, obj interface{}, opts *types.Opts) error {

	var (
		rev   *int64
		force bool
		ttl   uint64
	)

	if opts != nil {
		rev = opts.Rev
		force = opts.Force
		ttl = opts.Ttl
	}

	return s.client.store.Set(ctx, keyCreate(collection, name), obj, nil, ttl, force, rev)
}

func (s Storage) Del(ctx context.Context, collection string, name string) error {

	key := keyCreate(collection, name)

	if name == "" {
		key = keyCreate(collection)
	}

	return s.client.store.Del(ctx, key)
}

func (s Storage) Watch(ctx context.Context, collection string, event chan *types.WatcherEvent, opts *types.Opts) error {

	log.V(logLevel).Debug("%s:> watch %s", logPrefix, collection)

	const filter = `\b.+\/(.+)\b`

	client, destroy, err := s.getClient()
	if err != nil {
		log.V(logLevel).Errorf("%s:> watch err: %v", logPrefix, err)
		return err
	}
	defer destroy()

	var rev *int64
	if opts != nil {
		rev = opts.Rev
	}

	watcher, err := client.Watch(ctx, keyCreate(collection), "", rev)
	if err != nil {
		log.V(logLevel).Errorf("%s:> watch err: %v", logPrefix, err)
		return err
	}

	for {
		select {
		case <-ctx.Done():
			log.V(logLevel).Debugf("%s:> the user interrupted watch", logPrefix)
			watcher.Stop()
			return nil
		case res := <-watcher.ResultChan():

			if res == nil {
				continue
			}

			if res.Type == types.STORAGEERROREVENT {
				err := res.Object.(error)
				log.Errorf("%s:> watch err: %v", logPrefix, err)
				return err
			}

			r, _ := regexp.Compile(filter)
			keys := r.FindStringSubmatch(res.Key)
			if len(keys) == 0 {
				continue
			}

			e := new(types.WatcherEvent)
			e.Action = res.Type
			e.SelfLink = keys[1]
			e.Storage.Key = res.Key
			e.Storage.Revision = res.Rev

			match := strings.Split(res.Key, ":")

			if len(match) > 0 {
				e.Name = match[len(match)-1]
			} else {
				e.Name = keys[0]
			}

			if res.Type == types.STORAGEDELETEEVENT {
				e.Data = res.Object
				event <- e
				continue
			}

			e.Data = res.Object

			event <- e
		}
	}

	return nil
}

func (s Storage) Filter() types.Filter {
	return new(Filter)
}

func (s Storage) Key() types.Key {
	return new(Key)
}

func (s Storage) Collection() types.Collection {
	return new(Collection)
}

func (s Storage) getClient() (store.Store, store.DestroyFunc, error) {
	return s.client.store, s.client.dfunc, nil
}
