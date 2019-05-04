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

package controller

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/onedomain/lastbackend/pkg/controller/envs"
	"github.com/onedomain/lastbackend/pkg/controller/ipam"
	"github.com/onedomain/lastbackend/pkg/controller/runtime"
	l "github.com/onedomain/lastbackend/pkg/log"
	"github.com/onedomain/lastbackend/pkg/storage"
	"github.com/spf13/viper"
)

const defaultCIDR = "172.0.0.0/24"

// Daemon - controller entrypoint
func Daemon(v *viper.Viper) bool {

	var (
		env  = envs.Get()
		sigs = make(chan os.Signal)
		done = make(chan bool, 1)
	)

	log := l.New(v.GetInt("verbose"))

	log.Info("Start Controller")

	if !v.IsSet("storage") {
		log.Fatalf("Storage not configured")
	}

	stg, err := storage.Get(v)
	if err != nil {
		log.Fatalf("Cannot initialize storage: %s", err.Error())
	}
	env.SetStorage(stg)

	cidr := defaultCIDR
	if v.IsSet("service") && v.IsSet("service.cidr") {
		cidr = v.GetString("service.cidr")
	}

	ipm, err := ipam.New(cidr)
	if err != nil {
		log.Fatalf("Cannot initialize ipam service: %s", err.Error())
	}
	env.SetIPAM(ipm)

	// Initialize Runtime
	r := runtime.NewRuntime(context.Background())
	r.Loop()

	// Handle SIGINT and SIGTERM.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-sigs:
				done <- true
				return
			}
		}
	}()

	<-done

	log.Info("Handle SIGINT and SIGTERM.")
	return true
}
