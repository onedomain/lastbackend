//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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

package exporter

import (
	"github.com/onedomain/lastbackend/pkg/util/http"
	"github.com/onedomain/lastbackend/pkg/util/http/middleware"
)

var Routes = []http.Route{
	{Path: "/exporter", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Authenticate}, Handler: ExporterListH},
	{Path: "/exporter/{exporter}", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Authenticate}, Handler: ExporterInfoH},
	{Path: "/exporter/{exporter}", Method: http.MethodPut, Middleware: []http.Middleware{middleware.Authenticate}, Handler: ExporterConnectH},
	{Path: "/exporter/{exporter}", Method: http.MethodDelete, Middleware: []http.Middleware{middleware.Authenticate}, Handler: ExporterRemoveH},
	{Path: "/exporter/{exporter}/status", Method: http.MethodPut, Middleware: []http.Middleware{middleware.Authenticate}, Handler: ExporterSetStatusH},
}
