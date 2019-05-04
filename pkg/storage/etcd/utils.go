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

import "strings"

const keySeparator = "/"

func BuildServiceKey(namespace, name string) string {
	return keyCreate(namespace, name)
}

func BuildPodKey(namespace, service, deployment, name string) string {
	return keyCreate(namespace, service, deployment, name)
}

func BuildProcessKey(kind, hostname string) string {
	return keyCreate(kind, "process", hostname)
}

func BuildProcessLeadKey(kind string) string {
	return keyCreate(kind, "lead")
}

func BuildRouteKey(namespace string) string {
	return keyCreate(namespace)
}

func BuildSecretKey(namespace string) string {
	return keyCreate(namespace)
}

func BuildEndpointKey(namespace, name string) string {
	return keyCreate(namespace, name)
}

func BuildVolumeKey(namespace, name string) string {
	return keyCreate(namespace, name)
}

func keyCreate(val ...string) string {
	return strings.Join(val, keySeparator)
}
