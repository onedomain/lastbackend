//
// KULADO INC. CONFIDENTIAL
// __________________
//
// [2014] - [2018] KULADO INC.
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

package runtime

import (
	"github.com/onedomain/lastbackend/pkg/exporter/envs"
	"github.com/onedomain/lastbackend/pkg/exporter/logger"
	"github.com/onedomain/lastbackend/pkg/log"
	"os"
)

type Runtime struct {
	logger *logger.Logger
	port   uint16
	iface  string
}

type RuntimeOpts struct {
	Port   uint16
	Iface  string
	Logger *logger.LoggerOpts
}

func New(opts *RuntimeOpts) (r *Runtime, err error) {
	r = new(Runtime)
	r.port = opts.Port

	if opts.Logger != nil {
		lg, err := logger.New(opts.Logger)
		if err != nil {
			log.Errorf("can not init logger: %s", err.Error())
			os.Exit(1)
			return nil, err
		}
		envs.Get().SetLogger(lg)
	}
	return r, nil
}

func (r Runtime) Start() error {
	if r.logger != nil {
		if err := r.logger.Listen(); err != nil {
			log.Errorf("can not start logger listener: %s", err.Error())
			return err
		}
	}
	return nil
}
