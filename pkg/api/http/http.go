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

package http

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/onedomain/lastbackend/pkg/api/http/cluster"
	"github.com/onedomain/lastbackend/pkg/api/http/config"
	"github.com/onedomain/lastbackend/pkg/api/http/deployment"
	"github.com/onedomain/lastbackend/pkg/api/http/discovery"
	"github.com/onedomain/lastbackend/pkg/api/http/events"
	"github.com/onedomain/lastbackend/pkg/api/http/exporter"
	"github.com/onedomain/lastbackend/pkg/api/http/ingress"
	"github.com/onedomain/lastbackend/pkg/api/http/job"
	"github.com/onedomain/lastbackend/pkg/api/http/namespace"
	"github.com/onedomain/lastbackend/pkg/api/http/node"
	"github.com/onedomain/lastbackend/pkg/api/http/route"
	"github.com/onedomain/lastbackend/pkg/api/http/secret"
	"github.com/onedomain/lastbackend/pkg/api/http/service"
	"github.com/onedomain/lastbackend/pkg/api/http/task"
	"github.com/onedomain/lastbackend/pkg/api/http/volume"
	"github.com/onedomain/lastbackend/pkg/log"
	"github.com/onedomain/lastbackend/pkg/util/http"
	"github.com/onedomain/lastbackend/pkg/util/http/cors"
)

const (
	logLevel  = 2
	logPrefix = "api:http"
)

// Extends routes variable
var Routes = make([]http.Route, 0)

type HttpOpts struct {
	Insecure bool

	CertFile string
	KeyFile  string
	CaFile   string

	BearerToken string
}

func AddRoutes(r ...[]http.Route) {
	for i := range r {
		Routes = append(Routes, r[i]...)
	}
}

func init() {

	// Cluster
	AddRoutes(cluster.Routes)
	AddRoutes(node.Routes)
	AddRoutes(ingress.Routes)
	AddRoutes(exporter.Routes)
	AddRoutes(discovery.Routes)

	// Namespace
	AddRoutes(namespace.Routes)
	AddRoutes(secret.Routes)
	AddRoutes(config.Routes)
	AddRoutes(route.Routes)
	AddRoutes(service.Routes)
	AddRoutes(deployment.Routes)
	AddRoutes(volume.Routes)
	AddRoutes(ingress.Routes)
	AddRoutes(job.Routes)
	AddRoutes(task.Routes)

	// events
	AddRoutes(events.Routes)
}

func Listen(host string, port int, opts *HttpOpts) error {

	if opts == nil {
		opts = new(HttpOpts)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "access_token", opts.BearerToken)

	log.V(logLevel).Debugf("%s:> listen HTTP server on %s:%d", logPrefix, host, port)

	r := mux.NewRouter()
	r.Methods("OPTIONS").HandlerFunc(cors.Headers)

	var notFound http.MethodNotAllowedHandler
	r.NotFoundHandler = notFound

	var notAllowed http.MethodNotAllowedHandler
	r.MethodNotAllowedHandler = notAllowed

	for _, rt := range Routes {
		log.V(logLevel).Debugf("%s:> init route: %s", logPrefix, rt.Path)
		r.Handle(rt.Path, http.Handle(ctx, rt.Handler, rt.Middleware...)).Methods(rt.Method)
	}

	if opts.Insecure {
		log.V(logLevel).Debugf("%s:> run insecure http server on %d port", logPrefix, port)
		return http.Listen(host, port, r)
	}

	log.V(logLevel).Debugf("%s:> run http server with tls on %d port", logPrefix, port)
	return http.ListenWithTLS(host, port, opts.CaFile, opts.CertFile, opts.KeyFile, r)
}
