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

package network

import (
	"github.com/onedomain/lastbackend/pkg/distribution/errors"
	"strconv"
	"strings"
)

func ParsePortMap(s string) (uint16, string, error) {

	var (
		port  uint16
		proto string
		err   error
	)

	pm := strings.Split(s, "/")
	switch len(pm) {
	case 0:
		break
	case 1:
		p, err := strconv.ParseUint(pm[0], 10, 16)
		if err != nil {
			break
		}
		port = uint16(p)
		proto = "tcp"
		break
	case 2:
		p, err := strconv.ParseUint(pm[0], 10, 16)
		if err != nil {
			break
		}
		port = uint16(p)
		proto = strings.ToLower(pm[1])
		break
	default:
		err = errors.New("Invalid port map declaration")
		return port, proto, err
	}

	return port, proto, nil
}