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
	"github.com/onedomain/lastbackend/pkg/distribution/errors"
	"github.com/onedomain/lastbackend/pkg/util/validator"
	"io"
	"io/ioutil"
)

type JobRequest struct{}

func (JobRequest) Manifest() *JobManifest {
	return new(JobManifest)
}

func (j *JobManifest) Validate() *errors.Err {
	switch true {
	case j.Meta.Name != nil && !validator.IsJobName(*j.Meta.Name):
		return errors.New("job").BadParameter("name")
	case j.Meta.Description != nil && len(*j.Meta.Description) > DEFAULT_DESCRIPTION_LIMIT:
		return errors.New("job").BadParameter("description")
	case len(j.Spec.Task.Template.Containers) == 0:
		return errors.New("job").BadParameter("spec")
	case len(j.Spec.Task.Template.Containers) != 0:
		for _, container := range j.Spec.Task.Template.Containers {
			if len(container.Image.Name) == 0 {
				return errors.New("job").BadParameter("image")
			}
		}
	}

	return nil
}

func (j *JobManifest) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("job").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("job").Unknown(err)
	}

	err = json.Unmarshal(body, j)
	if err != nil {
		return errors.New("job").IncorrectJSON(err)
	}

	if err := j.Validate(); err != nil {
		return err
	}

	return nil
}

func (JobRequest) RemoveOptions() *JobRemoveOptions {
	return new(JobRemoveOptions)
}

func (s *JobRemoveOptions) Validate() *errors.Err {
	return nil
}
