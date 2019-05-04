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

const logTaskPrefix = "state:task:>"

type TaskState struct {
	lock  sync.RWMutex
	tasks map[string]types.NodeTask
}

func (s *TaskState) AddTask(key string, task *types.NodeTask) {
	log.V(logLevel).Debugf("%s add cancel func pod: %s", logTaskPrefix, key)
	s.tasks[key] = *task
}

func (s *TaskState) GetTask(key string) *types.NodeTask {
	log.V(logLevel).Debugf("%s get cancel func pod: %s", logTaskPrefix, key)

	if _, ok := s.tasks[key]; ok {
		t := s.tasks[key]
		return &t
	}

	return nil
}

func (s *TaskState) DelTask(pod *types.Pod) {
	log.V(logLevel).Debugf("%s del cancel func pod: %s", logTaskPrefix, pod.SelfLink())
	delete(s.tasks, pod.Meta.Name)
}
