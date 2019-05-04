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
	"io"
	"io/ioutil"
)

type TaskRequest struct{}

func (TaskRequest) Manifest() *TaskManifest {
	return new(TaskManifest)
}

func (t *TaskManifest) Validate() *errors.Err {
	switch true {
	}

	return nil
}

func (t *TaskManifest) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("service").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("task").Unknown(err)
	}

	err = json.Unmarshal(body, t)
	if err != nil {
		return errors.New("task").IncorrectJSON(err)
	}

	if err := t.Validate(); err != nil {
		return err
	}

	return nil
}

func (TaskRequest) CancelOptions() *TaskCancelOptions {
	return new(TaskCancelOptions)
}

func (s *TaskCancelOptions) Validate() *errors.Err {
	return nil
}

func (TaskRequest) LogOptions() *TaskLogsOptions {
	return new(TaskLogsOptions)
}

func (s *TaskLogsOptions) Validate() *errors.Err {
	return nil
}
