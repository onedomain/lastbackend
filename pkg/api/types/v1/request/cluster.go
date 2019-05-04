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

// ClusterUpdateOptions represents options availible to update in cluster
//
// swagger:model request_cluster_update
type ClusterUpdateOptions struct {
	Description *string               `json:"description"`
	Quotas      *ClusterQuotasOptions `json:"quotas"`
}

// swagger:model request_cluster_update_quotas
type ClusterQuotasOptions struct{}
