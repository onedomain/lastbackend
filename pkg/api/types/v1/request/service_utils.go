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
	"github.com/onedomain/lastbackend/pkg/util/validator"
)

type ServiceRequest struct{}

func (ServiceRequest) Manifest() *ServiceManifest {
	return new(ServiceManifest)
}

func (s *ServiceManifest) Validate() *errors.Err {
	switch true {
	case s.Meta.Name != nil && !validator.IsServiceName(*s.Meta.Name):
		return errors.New("service").BadParameter("name")
	case s.Meta.Description != nil && len(*s.Meta.Description) > DEFAULT_DESCRIPTION_LIMIT:
		return errors.New("service").BadParameter("description")
	case len(s.Spec.Template.Containers) == 0:
		return errors.New("service").BadParameter("spec")
	case len(s.Spec.Template.Containers) != 0:
		for _, container := range s.Spec.Template.Containers {
			if len(container.Image.Name) == 0 {
				return errors.New("service").BadParameter("image")
			}
		}
	}

	return nil
}

func (s *ServiceManifest) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("service").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("service").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("service").IncorrectJSON(err)
	}

	if err := s.Validate(); err != nil {
		return err
	}

	return nil
}

func (ServiceRequest) RemoveOptions() *ServiceRemoveOptions {
	return new(ServiceRemoveOptions)
}

func (s *ServiceRemoveOptions) Validate() *errors.Err {
	return nil
}
