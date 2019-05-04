//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package request

import (
	"encoding/json"
	"github.com/onedomain/lastbackend/pkg/distribution/types"
	"io"
	"io/ioutil"

	"github.com/onedomain/lastbackend/pkg/distribution/errors"
)

type VolumeRequest struct{}

func (VolumeRequest) Manifest() *VolumeManifest {
	return new(VolumeManifest)
}

func (v *VolumeManifest) Validate() *errors.Err {

	if v.Spec.Type == types.EmptyString {
		return errors.BadParameter("spec.type")
	}

	return nil
}

func (v *VolumeManifest) DecodeAndValidate(reader io.Reader) *errors.Err {

	if reader == nil {
		err := errors.New("data body can not be null")
		return errors.New("volume").IncorrectJSON(err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.New("volume").Unknown(err)
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return errors.New("volume").IncorrectJSON(err)
	}

	return v.Validate()
}

func (VolumeRequest) RemoveOptions() *VolumeRemoveOptions {
	return new(VolumeRemoveOptions)
}

func (v *VolumeRemoveOptions) Validate() *errors.Err {
	return nil
}
