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

// DeploymentUpdateOptions represents options availible to update in deployment
//
// swagger:model request_deployment_update
type DeploymentUpdateOptions struct {
	// Number of replicas
	// required: false
	Replicas *int `json:"replicas"`
	// Deployment status for update
	// required: false
	Status *struct {
		State   string `json:"state"`
		Message string `json:"message"`
	} `json:"status"`
}
