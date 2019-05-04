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

package exporter

import (
	"context"
	"github.com/onedomain/lastbackend/pkg/distribution/types"
	"github.com/onedomain/lastbackend/pkg/log"
	"sync"
)

const (
	logPrefix = "controller:>"
	logLevel  = 3
)

type Exporter struct {
	cache struct {
		lock        sync.RWMutex
		cluster     *types.Cluster
		nodes       map[string]*types.Node
		services    map[string]*types.Service
		deployments map[string]*types.Deployment
		pods        map[string]*types.Pod
		volumes     map[string]*types.Volume
		routes      map[string]*types.Route
	}
}

func New() *Exporter {
	var c = new(Exporter)
	return c
}

func (c *Exporter) Connect(ctx context.Context) error {
	log.V(logLevel).Debugf("%s:connect:> connect init", logPrefix)



	return nil
}



func (c *Exporter) SendClusterState(ctx context.Context) error {
	return nil
}

func (c *Exporter) SendNodeState(ctx context.Context) error {
	return nil
}

func (c *Exporter) SendServiceState(ctx context.Context) error {
	return nil
}

func (c *Exporter) SendDeploymentState(ctx context.Context) error {
	return nil
}

func (c *Exporter) SendPodState(ctx context.Context) error {
	return nil
}

func (c *Exporter) SendVolumeState(ctx context.Context) error {
	return nil
}
