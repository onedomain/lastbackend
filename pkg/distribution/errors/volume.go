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

package errors

import "errors"

const (
	VolumeIsNotReady                   = "volume is not ready"
	VolumeIsNotProvisioned             = "volume is not provisioned"
	VolumeNotFound                     = "volume not found"
	VolumesProvisionedOnDifferentNodes = "volumes are binded to different nodes"
)

type VolumeError struct {
}

func (ve *VolumeError) NotReady(vol string) error {
	return errors.New(joinNameAndMessage(vol, VolumeIsNotReady))
}

func (ve *VolumeError) NotProvisioned(vol string) error {
	return errors.New(joinNameAndMessage(vol, VolumeIsNotProvisioned))
}

func (ve *VolumeError) NotFound(vol string) error {
	return errors.New(joinNameAndMessage(vol, VolumeNotFound))
}

func (ve *VolumeError) DifferentNodes() error {
	return errors.New(VolumesProvisionedOnDifferentNodes)
}

func (e *err) Volume() *VolumeError {
	return new(VolumeError)
}
