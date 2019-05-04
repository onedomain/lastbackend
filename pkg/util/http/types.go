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
	"net/http"
)

type Route struct {
	Path       string
	Handler    func(w http.ResponseWriter, r *http.Request)
	Middleware []Middleware
	Method     string
}

//type Middleware interface {
//	Handler(func(w http.ResponseWriter, r *http.Request) error) func(w http.ResponseWriter, r *http.Request) error
//}

type Middleware func(context.Context, http.HandlerFunc) http.HandlerFunc
