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

type DeploymentRequest struct{}

func (DeploymentRequest) UpdateOptions() *DeploymentUpdateOptions {
	return new(DeploymentUpdateOptions)
}

func (d *DeploymentUpdateOptions) Validate() *errors.Err {
	switch true {
	case d.Replicas == nil:
		return errors.New("deployment").BadParameter("replicas")
	case *d.Replicas < 1:
		return errors.New("deployment").BadParameter("replicas")
	}
	return nil
}

func (d *DeploymentUpdateOptions) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("deployment").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("deployment").Unknown(err)
	}

	err = json.Unmarshal(body, d)
	if err != nil {
		return errors.New("deployment").IncorrectJSON(err)
	}

	return d.Validate()
}

func (s *DeploymentUpdateOptions) ToJson() ([]byte, error) {
	return json.Marshal(s)
}
