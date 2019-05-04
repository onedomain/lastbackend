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

// Last.Backend Open-source API
//
// Open-source system for automating deployment, scaling, and management of containerized applications.
//
// Terms Of Service:
//
// https://0xqi.com/legal/terms/
//
//     Schemes: https
//     Host: api.0xqi.com
//     BasePath: /
//     Version: 0.9.4
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Last.Backend Teams <team@0xqi.com> https://0xqi.com
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - bearerToken:
//
//     SecurityDefinitions:
//       bearerToken:
//         description: Bearer Token authentication
//         type: apiKey
//         name: authorization
//         in: header
//
//     Extensions:
//     x-meta-value: value
//     x-meta-array:
//       - value1
//       - value2
//     x-meta-array-obj:
//       - name: obj
//         value: field
//
// swagger:meta
package main

import (
	"fmt"
	"strings"

	"github.com/onedomain/lastbackend/pkg/exporter"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const default_env_prefix = "LB"
const default_config_type = "yaml"
const default_config_name = "config"

var (
	flags = []struct {
		// flag name
		Name string
		// flag short name
		Short string
		// flag value
		Value interface{}
		// flag description
		Desc string
		// viper name for binding from flag
		Bind string
	}{
		{Name: "access-token", Short: "", Value: "", Desc: "Access token to API server", Bind: "token"},
		{Name: "bind-listener-address", Short: "", Value: "0.0.0.0", Desc: "Exporter bind address", Bind: "logger.host"},
		{Name: "bind-listener-port", Short: "", Value: 2963, Desc: "Exporter bind port", Bind: "logger.port"},
		{Name: "bind-rest-address", Short: "", Value: "0.0.0.0", Desc: "Exporter REST listener address", Bind: "server.tls.host"},
		{Name: "bind-rest-port", Short: "", Value: 2964, Desc: "Exporter REST listener port", Bind: "server.tls.port"},
		{Name: "tls-cert-file", Short: "", Value: "", Desc: "Exporter REST TLS cert file path", Bind: "server.tls.cert"},
		{Name: "tls-private-key-file", Short: "", Value: "", Desc: "Exportter REST TLS private key path", Bind: "server.tls.key"},
		{Name: "tls-ca-file", Short: "", Value: "", Desc: "Exporter REST TLS certificate authority file path", Bind: "server.tls.ca"},
		{Name: "api-uri", Short: "", Value: "", Desc: "REST API endpoint", Bind: "api.uri"},
		{Name: "api-cert-file", Short: "", Value: "", Desc: "REST API TLS certificate file path", Bind: "api.tls.cert"},
		{Name: "api-private-key-file", Short: "", Value: "", Desc: "REST API TLS private key file path", Bind: "api.tls.key"},
		{Name: "api-ca-file", Short: "", Value: "", Desc: "REST API TSL certificate authority file path", Bind: "api.tls.ca"},
		{Name: "bind-interface", Short: "", Value: "eth0", Desc: "Exporter bind network interface", Bind: "network.interface"},
		{Name: "log-workdir", Short: "", Value: "/var/run/lastbackend", Desc: "Set directory on host for logs storage", Bind: "logger.workdir"},
		{Name: "verbose", Short: "v", Value: 0, Desc: "Set log level from 0 to 7", Bind: "verbose"},
		{Name: "config", Short: "c", Value: "", Desc: "Path for the configuration file", Bind: "config"},
	}
)

func main() {

	for _, item := range flags {
		switch item.Value.(type) {
		case string:
			flag.StringP(item.Name, item.Short, item.Value.(string), item.Desc)
		case int:
			flag.IntP(item.Name, item.Short, item.Value.(int), item.Desc)
		case []string:
			flag.StringSliceP(item.Name, item.Short, item.Value.([]string), item.Desc)
		default:
			panic(fmt.Sprintf("bad %s argument value", item.Name))
		}
	}

	flag.Parse()

	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.SetEnvPrefix(default_env_prefix)

	for _, item := range flags {
		if err := v.BindPFlag(item.Bind, flag.Lookup(item.Name)); err != nil {
			panic(err)
		}

		name := strings.Replace(strings.ToUpper(item.Name), "-", "_", -1)
		name = strings.Join([]string{default_env_prefix, name}, "_")

		if err := v.BindEnv(item.Bind, name); err != nil {
			panic(err)
		}
	}

	v.SetConfigType(default_config_type)
	v.SetConfigFile(v.GetString(default_config_name))

	if len(v.GetString("config")) != 0 {
		if err := v.ReadInConfig(); err != nil {
			panic(fmt.Sprintf("Read config err: %v", err))
		}
	}

	exporter.Daemon(v)
}
