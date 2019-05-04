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

package envs

import (
	"github.com/onedomain/lastbackend/pkg/controller/ipam/ipam"
	"github.com/onedomain/lastbackend/pkg/storage"
)

var e Env

type Env struct {
	storage storage.Storage
	ipam    ipam.IPAM
}

func Get() *Env {
	return &e
}

func (c *Env) SetStorage(storage storage.Storage) {
	c.storage = storage
}

func (c *Env) GetStorage() storage.Storage {
	return c.storage
}

func (c *Env) SetIPAM(ipam ipam.IPAM) {
	c.ipam = ipam
}

func (c *Env) GetIPAM() ipam.IPAM {
	return c.ipam
}
