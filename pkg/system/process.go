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

package system

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage"
	"github.com/lastbackend/lastbackend/pkg/util/system"

	"encoding/json"

	"github.com/lastbackend/dynamic/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/storage/etcd"
)

// HeartBeat Interval
const heartBeatInterval = 10 // in seconds
const systemLeadTTL = 11
const logLevel = 7

type Process struct {
	// Process storage
	storage storage.Storage
	// Managed process
	process *types.Process
}

// Process register function
// The main purpose is to register process in the system
// If we need to distribution and need master/replicas, use WaitElected function
func (c *Process) Register(ctx context.Context, kind string, stg storage.Storage) (*types.Process, error) {

	var (
		err  error
		item = new(types.Process)
	)

	log.V(logLevel).Debugf("System: Process: Register: %s", kind)
	item.Meta.SetDefault()
	item.Meta.Kind = kind

	if item.Meta.Hostname, err = system.GetHostname(); err != nil {
		log.Errorf("System: Process: Register: get hostname: %s", err.Error())
		return item, err
	}

	item.Meta.PID = system.GetPid()
	item.Meta.ID = encodeID(item)

	c.process = item
	c.storage = stg

	opts := storage.GetOpts()
	opts.Ttl = systemLeadTTL
	opts.Force = true
	if err := c.storage.Set(ctx, storage.SystemKind, c.storage.Key().Process(c.process.Meta.Hostname, false), c.process, opts); err != nil {
		log.Errorf("System: Process: Register: %s", err.Error())
		return item, err
	}

	go c.HeartBeat(ctx)

	return item, nil
}

// HeartBeat function - updates current process state
// and master election ttl option
func (c *Process) HeartBeat(ctx context.Context) {

	log.V(logLevel).Debugf("System: Process: Start HeartBeat for: %s", c.process.Meta.Kind)
	ticker := time.NewTicker(heartBeatInterval * time.Second)

	opts := storage.GetOpts()
	opts.Ttl = systemLeadTTL
	opts.Force = true

	for range ticker.C {
		// Update process state
		log.V(logLevel).Debug("System: Process: Beat")

		if err := c.storage.Set(ctx, storage.SystemKind, c.storage.Key().Process(c.process.Meta.Hostname, false), c.process, opts); err != nil {
			log.Errorf("System: Process: Register: %s", err.Error())
			return
		}

		// Check election
		if c.process.Meta.Lead {
			log.V(logLevel).Debug("System: Process: Beat: Lead TTL update")

			if err := c.storage.Set(ctx, storage.SystemKind, c.storage.Key().Process(c.process.Meta.Hostname, true), c.process, opts); err != nil {
				log.Errorf("System: Process: update process: %s", err.Error())
				return
			}

		}

	}
}

// WaitElected function used for election waiting if
// master/replicas type of process used
func (c *Process) WaitElected(ctx context.Context, lead chan bool) error {

	log.V(logLevel).Debug("System: Process: Wait for election")

	l := false
	opts := storage.GetOpts()
	opts.Ttl = systemLeadTTL

	if err := c.storage.Get(ctx, storage.SystemKind, etcd.BuildProcessLeadKey(c.process.Meta.Kind), &l); err != nil {
		log.Errorf("System: Process: get lead process: %s", err.Error())

		if err.Error() == store.ErrEntityNotFound {
			err = c.storage.Put(ctx, storage.SystemKind, c.storage.Key().Process(c.process.Meta.Hostname, true), c.process, opts)
			if err != nil && err.Error() != store.ErrEntityExists {
				log.V(logLevel).Errorf("System: Process: create process ttl err: %s", err.Error())
				return err
			}
			l = true
		} else {
			return err
		}

	}

	if l {
		log.V(logLevel).Debug("System: Process: Set as Lead")
		c.process.Meta.Lead = true
		c.process.Meta.Slave = false
		lead <- true
	}

	done := make(chan bool)
	watcher := storage.NewWatcher()

	go func() {
		for {
			select {
			case <-ctx.Done():
				done <- true
				return
			case e := <-watcher:
				if e.Data == nil {
					continue
				}

				switch e.Action {
				case types.EventActionCreate:
					fallthrough
				case types.EventActionUpdate:
					var id string
					if err := json.Unmarshal(e.Data.([]byte), &id); err != nil {
						continue
					}

					if id == c.process.Meta.ID {
						lead <- true
					} else {
						lead <- false
					}
				case types.EventActionDelete:

					if err := c.storage.Get(ctx, storage.SystemKind, etcd.BuildProcessLeadKey(c.process.Meta.Kind), &l); err != nil {
						log.Errorf("System: Process: get lead process: %s", err.Error())

						if err.Error() == store.ErrEntityNotFound {

							err = c.storage.Put(ctx, storage.SystemKind, c.storage.Key().Process(c.process.Meta.Hostname, true), c.process, opts)
							if err != nil && err.Error() != store.ErrEntityExists {
								log.V(logLevel).Errorf("System: Process: create process ttl err: %s", err.Error())
								continue
							}
							l = true
						} else {
							continue
						}

					}

				}

			}
		}
	}()

	if err := c.storage.Watch(ctx, storage.SystemKind, watcher); err != nil {
		return err
	}

	return nil
}

// Encode unique ID from pid and process hostname
func encodeID(c *types.Process) string {
	key := fmt.Sprintf("%s|%d", c.Meta.Hostname, c.Meta.PID)
	return base64.StdEncoding.EncodeToString([]byte(key))
}

// Decode ID into hostname and pid
func decodeID(id string) (int, string, error) {

	var (
		key      []byte
		pid      int
		err      error
		hostname string
	)

	key, err = base64.StdEncoding.DecodeString(id)
	if err != nil {
		return pid, hostname, err
	}

	parts := strings.Split(string(key), "|")
	if len(parts) == 2 {
		hostname = parts[0]
		pid, _ = strconv.Atoi(parts[1])
	}

	return pid, hostname, nil
}
