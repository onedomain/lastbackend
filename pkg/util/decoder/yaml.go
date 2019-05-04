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

package decoder

import (
	"bytes"
)

const separator = "\n---"

func YamlSplit(data []byte) [][]byte {

	sep := []byte(separator)
	obj := make([][]byte, 0)

	for {
		i := bytes.Index(data, sep)
		if i < 0 {
			break
		}

		item := make([]byte, i)
		if len(data) < i {
			break
		}

		copy(item, data[: i])
		obj = append(obj, item)
		data = data[i + len(sep):]
	}

	if len(data) > 0 {
		obj = append(obj, data)
	}

	return obj
}
