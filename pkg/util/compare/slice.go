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

package compare

func SliceOfString(a, b []string) bool {

	if len(a) != len(b) {
		return false
	}

	diff := make(map[string]int, len(a))

	for _, _a := range a {
		diff[_a]++
	}

	for _, _b := range b {

		if _, ok := diff[_b]; !ok {
			return false
		}

		diff[_b] -= 1
		if diff[_b] == 0 {
			delete(diff, _b)
		}

	}

	return len(diff) == 0
}
