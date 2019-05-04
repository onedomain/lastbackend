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

package errors

import "errors"

var (
	NotLoggedMessage   = errors.New("You are currently not logged in to the system, to get proper access create a new user or login with an existing user.")
	LoginErrorMessage  = errors.New("Incorrect login or password")
	LogoutErrorMessage = errors.New("Some problems with logout")
	UnknownMessage     = errors.New("Oops, error occurred: Please provide bug to github: https://github.com/onedomain/registry/issues")
)
