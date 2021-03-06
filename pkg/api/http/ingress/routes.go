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

package ingress

import (
	"github.com/onedomain/lastbackend/pkg/util/http"
	"github.com/onedomain/lastbackend/pkg/util/http/middleware"
)

var Routes = []http.Route{
	{Path: "/ingress", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Authenticate}, Handler: IngressListH},
	{Path: "/ingress/{ingress}", Method: http.MethodGet, Middleware: []http.Middleware{middleware.Authenticate}, Handler: IngressInfoH},
	{Path: "/ingress/{ingress}", Method: http.MethodPut, Middleware: []http.Middleware{middleware.Authenticate}, Handler: IngressConnectH},
	{Path: "/ingress/{ingress}", Method: http.MethodDelete, Middleware: []http.Middleware{middleware.Authenticate}, Handler: IngressRemoveH},
	{Path: "/ingress/{ingress}/status", Method: http.MethodPut, Middleware: []http.Middleware{middleware.Authenticate}, Handler: IngressSetStatusH},
}
