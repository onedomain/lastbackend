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

package state

import (
	"github.com/onedomain/lastbackend/pkg/distribution/types"
	"github.com/onedomain/lastbackend/pkg/log"
	"sync"
)

const logEndpointPrefix = "state:endpoints:>"

type EndpointState struct {
	lock      sync.RWMutex
	hash      string
	endpoints map[string]*types.EndpointState
}

func (es *EndpointState) GetHash() string {
	return es.hash
}

func (es *EndpointState) SetHash(hash string) {
	es.hash = hash
}

func (es *EndpointState) GetEndpoints() map[string]*types.EndpointState {
	log.V(logLevel).Debugf("%s get endpoints", logEndpointPrefix)
	return es.endpoints
}

func (es *EndpointState) SetEndpoints(endpoints map[string]*types.EndpointState) {
	es.lock.Lock()
	defer es.lock.Unlock()

	for key, endpoint := range endpoints {
		log.V(logLevel).Debugf("%s: add endpoint %s", logEndpointPrefix, key)
		es.endpoints[key] = endpoint
	}
}

func (es *EndpointState) GetEndpoint(key string) *types.EndpointState {
	log.V(logLevel).Debugf("%s: get endpoint: %s", logEndpointPrefix, key)
	es.lock.Lock()
	defer es.lock.Unlock()

	ep, ok := es.endpoints[key]
	if !ok {
		return nil
	}

	return ep
}

func (es *EndpointState) AddEndpoint(key string, endpoint *types.EndpointState) {
	log.V(logLevel).Debugf("%s: add endpoint %s", logEndpointPrefix, key)
	es.lock.Lock()
	defer es.lock.Unlock()
	es.endpoints[key] = endpoint
}

func (es *EndpointState) SetEndpoint(key string, endpoint *types.EndpointState) {
	log.V(logLevel).Debugf("%s: set endpoint %s", logEndpointPrefix, key)
	es.lock.Lock()
	defer es.lock.Unlock()
	es.endpoints[key] = endpoint
}

func (es *EndpointState) DelEndpoint(key string) {
	log.V(logLevel).Debugf("%s: del endpoint %s", logEndpointPrefix, key)
	es.lock.Lock()
	defer es.lock.Unlock()
	delete(es.endpoints, key)
}
