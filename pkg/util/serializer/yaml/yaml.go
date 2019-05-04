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

package yaml

import (
	"gopkg.in/yaml.v2"
	"io"
)

type Encoder struct{}
type Decoder struct{}

func (Encoder) Encode(objPtr interface{}, w io.Writer) error {
	buf, err := yaml.Marshal(objPtr)
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	return err
}

func (Decoder) Decode(data []byte, objPtr interface{}) error {
	return yaml.Unmarshal(data, objPtr)
}
