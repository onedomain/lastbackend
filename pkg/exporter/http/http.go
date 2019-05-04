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

package http

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/onedomain/lastbackend/pkg/exporter/http/logs"
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

	BearerToken string

	CertFile string
	KeyFile  string
	CaFile   string
}

func AddRoutes(r ...[]http.Route) {
	for i := range r {
		Routes = append(Routes, r[i]...)
	}
}

func init() {
	// Cluster
	AddRoutes(logs.Routes)
}

func Listen(host string, port int, opts *HttpOpts) error {

	if opts == nil {
		opts = new(HttpOpts)
	}

	log.V(logLevel).Debugf("%s:> listen HTTP server on %s:%d", logPrefix, host, port)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "access_token", opts.BearerToken)

	r := mux.NewRouter()
	r.Methods("OPTIONS").HandlerFunc(cors.Headers)

	var notFound http.MethodNotAllowedHandler
	r.NotFoundHandler = notFound

	var notAllowed http.MethodNotAllowedHandler
	r.MethodNotAllowedHandler = notAllowed

	for _, route := range Routes {
		log.V(logLevel).Debugf("%s:> init route: %s", logPrefix, route.Path)
		r.Handle(route.Path, http.Handle(ctx, route.Handler, route.Middleware...)).Methods(route.Method)
	}

	if opts.Insecure {
		log.V(logLevel).Debugf("%s:> run insecure http server", logPrefix)
		return http.Listen(host, port, r)
	}

	log.V(logLevel).Debugf("%s:> run http server with tls", logPrefix)
	return http.ListenWithTLS(host, port, opts.CaFile, opts.CertFile, opts.KeyFile, r)
}
