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

package distribution

import (
	"context"
	"github.com/onedomain/lastbackend/pkg/distribution/types"
	"github.com/onedomain/lastbackend/pkg/log"
	"github.com/onedomain/lastbackend/pkg/storage"

	"encoding/json"

	"github.com/onedomain/lastbackend/pkg/distribution/errors"
)

const (
	logClusterPrefix = "distribution:cluster"
)

// Cluster - distribution model
type Cluster struct {
	context context.Context
	storage storage.Storage
}

// Info - get cluster info
func (c *Cluster) Get() (*types.Cluster, error) {

	log.V(logLevel).Debugf("%s:get:> get info", logClusterPrefix)

	cluster := new(types.Cluster)
	err := c.storage.Get(c.context, c.storage.Collection().Cluster(), types.EmptyString, cluster, nil)
	if err != nil {
		if errors.Storage().IsErrEntityNotFound(err) {
			log.V(logLevel).Warnf("%s:get:> cluster not found", logClusterPrefix)
			return cluster, nil
		}

		log.V(logLevel).Errorf("%s:get:> get cluster err: %v", logClusterPrefix, err)
		return nil, err
	}

	return cluster, nil
}

func (c *Cluster) Set(cluster *types.Cluster) error {

	log.V(logLevel).Debugf("%s:set:> update Cluster %#v", logClusterPrefix, cluster)
	opts := storage.GetOpts()
	opts.Force = true

	if err := c.storage.Set(c.context, c.storage.Collection().Cluster(), types.EmptyString, cluster, opts); err != nil {
		log.V(logLevel).Errorf("%s:set:> update Cluster err: %v", logClusterPrefix, err)
		return err
	}

	return nil
}

// Watch cluster changes
func (c *Cluster) Watch(ch chan types.ClusterEvent) {

	log.V(logLevel).Debugf("%s:watch:> watch cluster", logClusterPrefix)

	done := make(chan bool)
	watcher := storage.NewWatcher()

	go func() {
		for {
			select {
			case <-c.context.Done():
				done <- true
				return
			case e := <-watcher:
				if e.Data == nil {
					continue
				}

				res := types.ClusterEvent{}
				res.Name = e.Name
				res.Action = e.Action

				cluster := new(types.Cluster)

				if err := json.Unmarshal(e.Data.([]byte), *cluster); err != nil {
					log.Errorf("%s:> parse data err: %v", logClusterPrefix, err)
					continue
				}

				res.Data = cluster

				ch <- res
			}
		}
	}()

	opts := storage.GetOpts()
	go c.storage.Watch(c.context, c.storage.Collection().Cluster(), watcher, opts)

	<-done
}

// NewClusterModel - return new cluster model
func NewClusterModel(ctx context.Context, stg storage.Storage) *Cluster {
	return &Cluster{ctx, stg}
}
