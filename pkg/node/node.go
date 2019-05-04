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

package node

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/onedomain/lastbackend/pkg/api/client"
	"github.com/onedomain/lastbackend/pkg/distribution/types"
	l "github.com/onedomain/lastbackend/pkg/log"
	"github.com/onedomain/lastbackend/pkg/network"
	"github.com/onedomain/lastbackend/pkg/node/controller"
	"github.com/onedomain/lastbackend/pkg/node/envs"
	"github.com/onedomain/lastbackend/pkg/node/exporter"
	"github.com/onedomain/lastbackend/pkg/node/http"
	"github.com/onedomain/lastbackend/pkg/node/runtime"
	"github.com/onedomain/lastbackend/pkg/node/state"
	"github.com/onedomain/lastbackend/pkg/runtime/cii/cii"
	"github.com/onedomain/lastbackend/pkg/runtime/cri/cri"
	"github.com/onedomain/lastbackend/pkg/runtime/csi/csi"
	"github.com/spf13/viper"
)

// Daemon - run node daemon
func Daemon(v *viper.Viper) {

	var (
		sigs = make(chan os.Signal)
		done = make(chan bool, 1)
	)

	log := l.New(v.GetInt("verbose"))
	log.Info("Start Node")

	if !v.IsSet("runtime") {
		log.Fatalf("Runtime not configured")
	}
	if !v.IsSet("runtime.cri") {
		log.Fatalf("CRI not configured")
	}
	if !v.IsSet("runtime.iri") {
		log.Fatalf("IRI not configured")
	}
	if !v.IsSet("runtime.csi") {
		log.Fatalf("CSI not configured")
	}

	if err := envs.Get().SetConfig(v); err != nil {
		log.Fatalf("Parse config err: %v", err)
	}

	criDriver := v.GetString("runtime.cri.type")
	_cri, err := cri.New(criDriver, v.GetStringMap(fmt.Sprintf("runtime.%s", criDriver)))
	if err != nil {
		log.Errorf("Cannot initialize cri: %v", err)
	}

	iriDriver := v.GetString("runtime.iri.type")
	_cii, err := cii.New(iriDriver, v.GetStringMap(fmt.Sprintf("runtime.%s", iriDriver)))
	if err != nil {
		log.Errorf("Cannot initialize iri: %v", err)
	}

	csis := v.GetStringMap("runtime.csi")
	if csis != nil {
		for kind := range csis {
			si, err := csi.New(kind, v)
			if err != nil {
				log.Errorf("Cannot initialize csi: %s > %v", kind, err)
				return
			}
			envs.Get().SetCSI(kind, si)
		}
	}

	st := state.New()
	envs.Get().SetState(st)
	envs.Get().SetCRI(_cri)
	envs.Get().SetCII(_cii)

	if v.IsSet("network") {
		net, err := network.New(v)
		if err != nil {
			log.Errorf("can not initialize network: %s", err.Error())
			os.Exit(1)
		}
		envs.Get().SetNet(net)
	}

	r, err := runtime.New()
	if err != nil {
		log.Errorf("can not initialize runtime: %s", err.Error())
		os.Exit(1)
	}

	st.Node().Info = runtime.NodeInfo()
	st.Node().Status = runtime.NodeStatus()

	if err := r.Restore(); err != nil {
		log.Errorf("restore err: %v", err)
		os.Exit(1)
	}
	r.Subscribe()
	r.Loop()

	c, err := exporter.NewExporter(st.Node().Info.Hostname, types.EmptyString)
	if err != nil {
		log.Errorf("can not initialize collector: %s", err.Error())
		os.Exit(1)
	}
	envs.Get().SetExporter(c)
	go c.Listen()

	if v.IsSet("manifest.dir") {
		dir := v.GetString("manifest.dir")
		if dir != types.EmptyString {
			r.Provision(dir)
		}
	}

	if v.IsSet("api") {

		cfg := client.NewConfig()
		cfg.BearerToken = v.GetString("token")

		if v.IsSet("api.tls") && !v.GetBool("api.tls.insecure") {
			cfg.TLS = client.NewTLSConfig()
			cfg.TLS.CertFile = v.GetString("api.tls.cert")
			cfg.TLS.KeyFile = v.GetString("api.tls.key")
			cfg.TLS.CAFile = v.GetString("api.tls.ca")
		}

		endpoint := v.GetString("api.uri")

		rest, err := client.New(client.ClientHTTP, endpoint, cfg)
		if err != nil {
			log.Errorf("Init client err: %s", err)
		}

		n := rest.V1().Cluster().Node(st.Node().Info.Hostname)
		s := rest.V1()
		envs.Get().SetClient(n, s)

		ctl := controller.New(r)

		if err := ctl.Connect(v); err != nil {
			log.Errorf("node:initialize: connect err %s", err.Error())
		}
		go ctl.Subscribe()
		go ctl.Sync()
	}

	go func() {
		opts := new(http.HttpOpts)
		opts.BearerToken = v.GetString("token")
		opts.Insecure = v.GetBool("server.tls.insecure")
		opts.CertFile = v.GetString("server.tls.server_cert")
		opts.KeyFile = v.GetString("server.tls.server_key")
		opts.CaFile = v.GetString("server.tls.ca")

		if err := http.Listen(v.GetString("server.host"), v.GetInt("server.port"), opts); err != nil {
			log.Fatalf("Http server start error: %v", err)
		}
	}()

	// Handle SIGINT and SIGTERM.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-sigs:
				r.Stop()
				done <- true
				return
			}
		}
	}()

	<-done
	log.Info("Handle SIGINT and SIGTERM.")

	return
}
