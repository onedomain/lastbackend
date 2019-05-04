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

package cache

import (
	"sync"
	"time"
)

type Item struct {
	sync.RWMutex
	data    []string
	expires *time.Time
}

func (item *Item) setExpireTime(duration time.Duration) {
	item.Lock()
	defer item.Unlock()

	expires := time.Now().Add(duration)
	item.expires = &expires
}

func (item *Item) expired() bool {
	item.RLock()
	defer item.RUnlock()
	return item.expires == nil || item.expires.Before(time.Now())
}
