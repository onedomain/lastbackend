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

package request

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/onedomain/lastbackend/pkg/distribution/errors"
)

type ClusterRequest struct{}

func (ClusterRequest) UpdateOptions() *ClusterUpdateOptions {
	return new(ClusterUpdateOptions)
}

func (c *ClusterUpdateOptions) Validate() *errors.Err {
	switch true {
	case c.Description != nil && len(*c.Description) > DEFAULT_DESCRIPTION_LIMIT:
		return errors.New("cluster").BadParameter("description")
	}
	return nil
}

func (c *ClusterUpdateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("cluster").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("cluster").Unknown(err)
	}

	err = json.Unmarshal(body, c)
	if err != nil {
		return errors.New("cluster").IncorrectJSON(err)
	}

	return c.Validate()
}

func (s *ClusterUpdateOptions) ToJson() ([]byte, error) {
	return json.Marshal(s)
}
