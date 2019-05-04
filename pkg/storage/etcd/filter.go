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

package etcd

import (
	"fmt"

	"github.com/onedomain/lastbackend/pkg/storage/types"
)

type Filter struct{}

func (Filter) Namespace() types.NamespaceFilter {
	return new(NamespaceFilter)
}

func (Filter) Service() types.ServiceFilter {
	return new(ServiceFilter)
}

func (Filter) Deployment() types.DeploymentFilter {
	return new(DeploymentFilter)
}

func (Filter) Pod() types.PodFilter {
	return new(PodFilter)
}

func (Filter) Endpoint() types.EndpointFilter {
	return new(EndpointFilter)
}

func (Filter) Route() types.RouteFilter {
	return new(RouteFilter)
}

func (Filter) Config() types.ConfigFilter {
	return new(ConfigFilter)
}

func (Filter) Secret() types.SecretFilter {
	return new(SecretFilter)
}

func (Filter) Volume() types.VolumeFilter {
	return new(VolumeFilter)
}

func (Filter) Task() types.TaskFilter {
	return new(TaskFilter)
}

func (Filter) Job() types.JobFilter {
	return new(JobFilter)
}

type NamespaceFilter struct{}

type ServiceFilter struct{}

func byNamespace(namespace string) string {
	return fmt.Sprintf("%s:", namespace)
}

func byService(namespace, service string) string {
	return fmt.Sprintf("%s:%s:", namespace, service)
}

func byDeployment(namespace, service, deployment string) string {
	return fmt.Sprintf("%s:%s:d_%s:", namespace, service, deployment)
}

func byJob(namespace, job string) string {

	if job == "" {
		job = "manual"
	}

	return fmt.Sprintf("%s:%s:", namespace, job)
}

func (ServiceFilter) ByNamespace(namespace string) string {
	return byNamespace(namespace)
}

type DeploymentFilter struct{}

func (DeploymentFilter) ByNamespace(namespace string) string {
	return byNamespace(namespace)
}

func (DeploymentFilter) ByService(namespace, service string) string {
	return byService(namespace, service)
}

type PodFilter struct{}

func (PodFilter) ByNamespace(namespace string) string {
	return byNamespace(namespace)
}

func (PodFilter) ByService(namespace, service string) string {
	return byService(namespace, service)
}

func (PodFilter) ByDeployment(namespace, service, deployment string) string {
	return byDeployment(namespace, service, deployment)
}

func (PodFilter) ByJob(namespace, job string) string {
	return byService(namespace, job)
}

func (PodFilter) ByTask(namespace, job, task string) string {
	return byDeployment(namespace, job, task)
}

type EndpointFilter struct{}

func (EndpointFilter) ByNamespace(namespace string) string {
	return byNamespace(namespace)
}

type RouteFilter struct{}

func (RouteFilter) ByNamespace(namespace string) string {
	return byNamespace(namespace)
}

type SecretFilter struct{}

func (SecretFilter) ByNamespace(namespace string) string {
	return byNamespace(namespace)
}

type ConfigFilter struct{}

func (ConfigFilter) ByNamespace(namespace string) string {
	return byNamespace(namespace)
}

type VolumeFilter struct{}

func (VolumeFilter) ByNamespace(namespace string) string {
	return byNamespace(namespace)
}

type ManifestFilter struct{}

func (ManifestFilter) ByNodeManifest(node string) string {
	return fmt.Sprintf("%s/", node)
}

func (ManifestFilter) ByKindManifest(node string, kind types.Kind) string {
	return fmt.Sprintf("%s/%s/", node, kind)
}

type TaskFilter struct{}

func (TaskFilter) ByNamespace(namespace string) string {
	return byNamespace(namespace)
}

func (TaskFilter) ByJob(namespace, job string) string {
	return byJob(namespace, job)
}

type JobFilter struct{}

func (JobFilter) ByNamespace(namespace string) string {
	return byNamespace(namespace)
}
