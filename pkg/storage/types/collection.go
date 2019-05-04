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

package types

type Collection interface {
	Namespace() string
	Service() string
	Deployment() string
	Cluster() string
	Pod() string
	Ingress() IngressCollection
	Exporter() ExporterCollection
	Discovery() DiscoveryCollection
	System() string
	Node() NodeCollection
	Route() string
	Volume() string
	Secret() string
	Config() string
	Endpoint() string
	Network() string
	Subnet() string
	Manifest() ManifestCollection
	Job() string
	Task() string
	Test() string
	Root() string
}

type ManifestCollection interface {
	Node() string
	Cluster() string
	Pod(node string) string
	Volume(node string) string
	Route(ingress string) string
	Ingress() string
	Subnet() string
	Secret() string
	Endpoint() string
}

type NodeCollection interface {
	Info() string
	Status() string
}

type DiscoveryCollection interface {
	Info() string
	Status() string
}

type ExporterCollection interface {
	Info() string
	Status() string
}

type IngressCollection interface {
	Info() string
	Status() string
}
