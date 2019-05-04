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

package service_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/onedomain/lastbackend/pkg/util/resource"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/onedomain/lastbackend/pkg/api/envs"
	"github.com/onedomain/lastbackend/pkg/api/http/service"
	"github.com/onedomain/lastbackend/pkg/api/types/v1"
	"github.com/onedomain/lastbackend/pkg/api/types/v1/request"
	"github.com/onedomain/lastbackend/pkg/api/types/v1/views"
	"github.com/onedomain/lastbackend/pkg/distribution/errors"
	"github.com/onedomain/lastbackend/pkg/distribution/types"
	"github.com/onedomain/lastbackend/pkg/storage"
	"github.com/stretchr/testify/assert"
)

// Testing ServiceInfoH handler
func TestServiceInfo(t *testing.T) {

	var ctx = context.Background()

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")
	s1 := getServiceAsset(ns1.Meta.Name, "demo", "")
	s2 := getServiceAsset(ns1.Meta.Name, "test", "")

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
		service   *types.Service
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		want         *views.Service
		wantErr      bool
		err          string
		expectedCode int
	}{
		{
			name:         "checking get service if not exists",
			handler:      service.ServiceInfoH,
			args:         args{ctx, ns1, s2},
			fields:       fields{stg},
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Service not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get service if namespace not exists",
			handler:      service.ServiceInfoH,
			args:         args{ctx, ns2, s1},
			fields:       fields{stg},
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get service successfully",
			handler:      service.ServiceInfoH,
			args:         args{ctx, ns1, s1},
			fields:       fields{stg},
			want:         v1.View().Service().NewWithDeployment(s1, nil, nil),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Service(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Service(), s1.SelfLink().String(), s1, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("GET", fmt.Sprintf("/namespace/%s/service/%s", tc.args.namespace.Meta.Name, tc.args.service.Meta.Name), nil)
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/namespace/{namespace}/service/{service}", tc.handler)

			setRequestVars(r, req)

			// We create assert ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			res := httptest.NewRecorder()

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			r.ServeHTTP(res, req)

			// Check the status code is what we expect.
			if !assert.Equal(t, tc.expectedCode, res.Code, "status code not equal") {
				return
			}

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			if tc.wantErr {
				assert.Equal(t, tc.err, string(body), "incorrect status code")
			} else {

				s := new(views.Service)
				err := json.Unmarshal(body, &s)
				assert.NoError(t, err)

				assert.Equal(t, tc.want.Meta.Name, s.Meta.Name, "name not equal")
			}

		})
	}

}

// Testing ServiceListH handler
func TestServiceList(t *testing.T) {

	var ctx = context.Background()

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")
	s1 := getServiceAsset(ns1.Meta.Name, "demo", "")
	s2 := getServiceAsset(ns1.Meta.Name, "test", "")

	sl := types.NewServiceMap()
	sl.Items[s1.SelfLink().String()] = s1
	sl.Items[s2.SelfLink().String()] = s2

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
		service   *types.Service
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		err          string
		want         *types.ServiceMap
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking get services list if namespace not found",
			args:         args{ctx, ns2, nil},
			fields:       fields{stg},
			handler:      service.ServiceListH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get services list successfully",
			args:         args{ctx, ns1, nil},
			fields:       fields{stg},
			handler:      service.ServiceListH,
			want:         sl,
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Service(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Service(), s1.SelfLink().String(), s1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Service(), s2.SelfLink().String(), s2, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("GET", fmt.Sprintf("/namespace/%s", tc.args.namespace.Meta.Name), nil)
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/namespace/{namespace}", tc.handler)

			setRequestVars(r, req)

			// We create assert ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			res := httptest.NewRecorder()

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			r.ServeHTTP(res, req)

			// Check the status code is what we expect.
			if !assert.Equal(t, tc.expectedCode, res.Code, "status code not equal") {
				return
			}

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			if tc.wantErr {
				assert.Equal(t, tc.err, string(body), "incorrect status code")
			} else {

				s := new(views.RouteList)
				err := json.Unmarshal(body, &s)
				assert.NoError(t, err)

				for _, item := range *s {
					if _, ok := tc.want.Items[item.Meta.SelfLink]; !ok {
						assert.Error(t, errors.New("not equals"))
					}
				}
			}

		})
	}

}

// Testing ServiceCreateH handler
func TestServiceCreate(t *testing.T) {

	var ctx = context.Background()

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")

	ns3 := getNamespaceAsset("limits", "")
	ns3.Spec.Resources.Limits.RAM, _ = resource.DecodeMemoryResource("1GB")
	ns3.Spec.Resources.Limits.CPU, _ = resource.DecodeCpuResource("1")

	s1 := getServiceAsset(ns1.Meta.Name, "demo", "")
	s2 := getServiceAsset(ns1.Meta.Name, "success", "")
	s3 := getServiceAsset(ns3.Meta.Name, "success", "")

	sm1 := getServiceManifest("errored", "image")
	sm1.Spec.Template.Containers[0].Resources = new(request.ManifestSpecTemplateContainerResources)
	sm1.Spec.Template.Containers[0].Resources.Limits = new(request.ManifestSpecTemplateContainerResource)
	sm1.Spec.Template.Containers[0].Resources.Limits.RAM = "0.5GB"

	sm2 := getServiceManifest("errored", "image")
	sm2.Spec.Template.Containers[0].Resources = new(request.ManifestSpecTemplateContainerResources)
	sm2.Spec.Template.Containers[0].Resources.Limits = new(request.ManifestSpecTemplateContainerResource)
	sm2.Spec.Template.Containers[0].Resources.Limits.RAM = "2GB"
	sm2.Spec.Template.Containers[0].Resources.Limits.CPU = "0.5"

	sm3 := getServiceManifest("errored", "image")
	sm3.Spec.Template.Containers[0].Resources = new(request.ManifestSpecTemplateContainerResources)
	sm3.Spec.Template.Containers[0].Resources.Limits = new(request.ManifestSpecTemplateContainerResource)
	sm3.Spec.Template.Containers[0].Resources.Limits.RAM = "512MB"
	sm3.Spec.Template.Containers[0].Resources.Limits.CPU = "1.5"

	sm4 := getServiceManifest("errored", "image")
	sm4.Spec.Template.Containers[0].Resources = new(request.ManifestSpecTemplateContainerResources)
	sm4.Spec.Template.Containers[0].Resources.Limits = new(request.ManifestSpecTemplateContainerResource)
	sm4.Spec.Template.Containers[0].Resources.Limits.RAM = "2GB"
	sm4.Spec.Template.Containers[0].Resources.Limits.CPU = "1.5"

	var rsm5 = 3
	sm5 := getServiceManifest("errored", "image")
	sm5.Spec.Template.Containers[0].Resources = new(request.ManifestSpecTemplateContainerResources)
	sm5.Spec.Template.Containers[0].Resources.Limits = new(request.ManifestSpecTemplateContainerResource)
	sm5.Spec.Replicas = &rsm5
	sm5.Spec.Template.Containers[0].Resources.Limits.RAM = "128MB"
	sm5.Spec.Template.Containers[0].Resources.Limits.CPU = "0.5"

	sm6 := getServiceManifest("success", "image")
	sm6.Spec.Template.Containers[0].Resources = new(request.ManifestSpecTemplateContainerResources)
	sm6.Spec.Template.Containers[0].Resources.Limits = new(request.ManifestSpecTemplateContainerResource)
	sm6.Spec.Template.Containers[0].Resources.Limits.RAM = "512MB"
	sm6.Spec.Template.Containers[0].Resources.Limits.CPU = "0.5"

	sm7 := getServiceManifest("success", "image")
	sm7.Spec.Template.Containers[0].Resources = new(request.ManifestSpecTemplateContainerResources)
	sm7.Spec.Template.Containers[0].Resources.Limits = new(request.ManifestSpecTemplateContainerResource)
	sm7.Spec.Template.Containers[0].Resources.Limits.RAM = ""
	sm7.Spec.Template.Containers[0].Resources.Limits.CPU = ""

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
		service   *types.Service
	}

	tests := []struct {
		name         string
		args         args
		fields       fields
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		data         *request.ServiceManifest
		want         *views.Service
		wantErr      bool
		err          string
		expectedCode int
	}{
		{
			name:         "checking create service if name already exists",
			args:         args{ctx, ns1, s1},
			fields:       fields{stg},
			handler:      service.ServiceCreateH,
			data:         getServiceManifest("demo", "redis"),
			err:          "{\"code\":400,\"status\":\"Not Unique\",\"message\":\"Name is already in use\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "checking create service if namespace not found",
			args:         args{ctx, ns2, s2},
			fields:       fields{stg},
			handler:      service.ServiceCreateH,
			data:         getServiceManifest("test", "redis"),
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking create service with replica 1 with ram limit are higher then namespace limits",
			args:         args{ctx, ns3, s2},
			fields:       fields{stg},
			handler:      service.ServiceCreateH,
			data:         sm2,
			err:          "{\"code\":400,\"status\":\"Bad Request\",\"message\":\"resources ram limit exceeded\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "checking create service with replica 1 with cpu limit are higher then namespace limits",
			args:         args{ctx, ns3, s2},
			fields:       fields{stg},
			handler:      service.ServiceCreateH,
			data:         sm3,
			err:          "{\"code\":400,\"status\":\"Bad Request\",\"message\":\"resources cpu limit exceeded\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "checking create service with replica 1 with ram & cpu limits are higher then namespace limits",
			args:         args{ctx, ns3, s2},
			fields:       fields{stg},
			handler:      service.ServiceCreateH,
			data:         sm4,
			err:          "{\"code\":400,\"status\":\"Bad Request\",\"message\":\"resources ram limit exceeded\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "checking create service with replica 3 with limits are higher then namespace limits",
			args:         args{ctx, ns3, s2},
			fields:       fields{stg},
			handler:      service.ServiceCreateH,
			data:         sm5,
			err:          "{\"code\":400,\"status\":\"Bad Request\",\"message\":\"resources cpu limit exceeded\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "check create service if bad parameter name",
			args:         args{ctx, ns1, s3},
			fields:       fields{stg},
			handler:      service.ServiceCreateH,
			data:         getServiceManifest("_____test", "redis"),
			err:          "{\"code\":400,\"status\":\"Bad Parameter\",\"message\":\"Bad name parameter\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		// TODO: check another spec parameters
		{
			name:         "check create service success",
			args:         args{ctx, ns1, s2},
			fields:       fields{stg},
			handler:      service.ServiceCreateH,
			data:         getServiceManifest("success", "redis"),
			want:         v1.View().Service().NewWithDeployment(s2, nil, nil),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "check create service success with default limits",
			args:         args{ctx, ns3, s3},
			fields:       fields{stg},
			handler:      service.ServiceCreateH,
			data:         sm7,
			want:         v1.View().Service().NewWithDeployment(s3, nil, nil),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "check create service success with limits",
			args:         args{ctx, ns3, s3},
			fields:       fields{stg},
			handler:      service.ServiceCreateH,
			data:         sm6,
			want:         v1.View().Service().NewWithDeployment(s3, nil, nil),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Service(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns3.SelfLink().String(), ns3, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Service(), s1.SelfLink().String(), s1, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			bd, err := tc.data.ToJson()
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", fmt.Sprintf("/namespace/%s/service", tc.args.namespace.Meta.Name), strings.NewReader(string(bd)))
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/namespace/{namespace}/service", tc.handler)

			setRequestVars(r, req)

			// We create assert ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			res := httptest.NewRecorder()

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			r.ServeHTTP(res, req)

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)
			// Check the status code is what we expect.
			if !assert.Equal(t, tc.expectedCode, res.Code, "status code not equal") {
				return
			}

			if tc.wantErr {
				assert.Equal(t, tc.err, string(body), "incorrect status code")
			} else {

				got := new(types.Service)
				err := tc.fields.stg.Get(tc.args.ctx, stg.Collection().Service(), tc.args.service.SelfLink().String(), got, nil)
				assert.NoError(t, err)

				if got == nil {
					t.Error("can not be not nil")
					return
				}

				assert.Equal(t, s3.Meta.Name, got.Meta.Name, "name not equal")
				assert.Equal(t, s3.Meta.Description, got.Meta.Description, "description not equal")
			}
		})
	}

}

// Testing ServiceUpdateH handler
func TestServiceUpdate(t *testing.T) {

	var ctx = context.Background()

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")
	ns3 := getNamespaceAsset("limits", "")

	ns3.Status.Resources.Allocated.RAM, _ = resource.DecodeMemoryResource("1.5GB")
	ns3.Status.Resources.Allocated.CPU, _ = resource.DecodeCpuResource("1.5")

	ns3.Spec.Resources.Limits.RAM, _ = resource.DecodeMemoryResource("2GB")
	ns3.Spec.Resources.Limits.CPU, _ = resource.DecodeCpuResource("2")

	s1 := getServiceAsset(ns1.Meta.Name, "demo", "")
	s2 := getServiceAsset(ns1.Meta.Name, "test", "")
	s3 := getServiceAsset(ns1.Meta.Name, "demo", "demo description")

	s4 := getServiceAsset(ns3.Meta.Name, "limited", "demo description")
	s4.Spec.Template.Containers[0].Resources.Limits.RAM, _ = resource.DecodeMemoryResource("512MB")
	s4.Spec.Template.Containers[0].Resources.Limits.CPU, _ = resource.DecodeCpuResource("0.5")

	m1 := getServiceManifest(s3.Meta.Name, "redis")
	m1.SetServiceSpec(s1)

	m3 := getServiceManifest(s3.Meta.Name, "redis")

	m3.Meta.Description = &s3.Meta.Description
	m3.Spec.Template.Containers[0].Env[0].Name = "updated"
	m3.Spec.Template.Containers[0].Env[1].Value = "meta"
	m3.Spec.Template.Volumes[0].Name = "secret-test"
	m3.Spec.Template.Volumes[0].Secret.Name = "r"

	m3.SetServiceSpec(s3)

	sm1 := getServiceManifest("limited", "image")
	sm1.Spec.Template.Containers[0].Resources = new(request.ManifestSpecTemplateContainerResources)
	sm1.Spec.Template.Containers[0].Resources.Limits = new(request.ManifestSpecTemplateContainerResource)
	sm1.Spec.Template.Containers[0].Resources.Limits.RAM = "0.5GB"

	sm2 := getServiceManifest("limited", "image")
	sm2.Spec.Template.Containers[0].Resources = new(request.ManifestSpecTemplateContainerResources)
	sm2.Spec.Template.Containers[0].Resources.Limits = new(request.ManifestSpecTemplateContainerResource)
	sm2.Spec.Template.Containers[0].Resources.Limits.RAM = "2GB"
	sm2.Spec.Template.Containers[0].Resources.Limits.CPU = "0.5"

	sm3 := getServiceManifest("limited", "image")
	sm3.Spec.Template.Containers[0].Resources = new(request.ManifestSpecTemplateContainerResources)
	sm3.Spec.Template.Containers[0].Resources.Limits = new(request.ManifestSpecTemplateContainerResource)
	sm3.Spec.Template.Containers[0].Resources.Limits.RAM = "512MB"
	sm3.Spec.Template.Containers[0].Resources.Limits.CPU = "1.5"

	sm4 := getServiceManifest("limited", "image")
	sm4.Spec.Template.Containers[0].Resources = new(request.ManifestSpecTemplateContainerResources)
	sm4.Spec.Template.Containers[0].Resources.Limits = new(request.ManifestSpecTemplateContainerResource)
	sm4.Spec.Template.Containers[0].Resources.Limits.RAM = "2GB"
	sm4.Spec.Template.Containers[0].Resources.Limits.CPU = "1.5"

	var rsm5 = 3
	sm5 := getServiceManifest("limited", "image")
	sm5.Spec.Replicas = &rsm5
	sm5.Spec.Template.Containers[0].Resources = new(request.ManifestSpecTemplateContainerResources)
	sm5.Spec.Template.Containers[0].Resources.Limits = new(request.ManifestSpecTemplateContainerResource)
	sm5.Spec.Template.Containers[0].Resources.Limits.RAM = "128MB"
	sm5.Spec.Template.Containers[0].Resources.Limits.CPU = "0.5"

	sm6 := getServiceManifest("limited", "image")
	sm6.Spec.Template.Containers[0].Resources = new(request.ManifestSpecTemplateContainerResources)
	sm6.Spec.Template.Containers[0].Resources.Limits = new(request.ManifestSpecTemplateContainerResource)
	sm6.Spec.Template.Containers[0].Resources.Limits.RAM = "600MB"
	sm6.Spec.Template.Containers[0].Resources.Limits.CPU = "0.6"

	sm7 := getServiceManifest("limited", "image")
	sm7.Spec.Template.Containers[0].Resources = new(request.ManifestSpecTemplateContainerResources)
	sm7.Spec.Template.Containers[0].Resources.Limits = new(request.ManifestSpecTemplateContainerResource)
	sm7.Spec.Template.Containers[0].Resources.Limits.RAM = "512MB"
	sm7.Spec.Template.Containers[0].Resources.Limits.CPU = "0.5"

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
		service   *types.Service
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		data         *request.ServiceManifest
		want         *views.Service
		wantErr      bool
		err          string
		expectedCode int
	}{
		{
			name:         "checking update service if name not exists",
			fields:       fields{stg},
			args:         args{ctx, ns1, s2},
			handler:      service.ServiceUpdateH,
			data:         getServiceManifest("test", "redis"),
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Service not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking update service if namespace not found",
			fields:       fields{stg},
			args:         args{ctx, ns2, s1},
			handler:      service.ServiceUpdateH,
			data:         getServiceManifest("demo", "redis"),
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking update service without limits if namespace with limits",
			args:         args{ctx, ns3, s4},
			fields:       fields{stg},
			handler:      service.ServiceUpdateH,
			data:         sm1,
			err:          "{\"code\":400,\"status\":\"Bad Request\",\"message\":\"resources cpu limit is required\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "checking update service with replica 1 with ram limit are higher then namespace limits",
			args:         args{ctx, ns3, s4},
			fields:       fields{stg},
			handler:      service.ServiceUpdateH,
			data:         sm2,
			err:          "{\"code\":400,\"status\":\"Bad Request\",\"message\":\"resources ram limit exceeded\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "checking update service with replica 1 with cpu limit are higher then namespace limits",
			args:         args{ctx, ns3, s4},
			fields:       fields{stg},
			handler:      service.ServiceUpdateH,
			data:         sm3,
			err:          "{\"code\":400,\"status\":\"Bad Request\",\"message\":\"resources cpu limit exceeded\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "checking update service with replica 1 with ram & cpu limits are higher then namespace limits",
			args:         args{ctx, ns3, s4},
			fields:       fields{stg},
			handler:      service.ServiceUpdateH,
			data:         sm4,
			err:          "{\"code\":400,\"status\":\"Bad Request\",\"message\":\"resources ram limit exceeded\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "checking update service with replica 3 with limits are higher then namespace limits",
			args:         args{ctx, ns3, s4},
			fields:       fields{stg},
			handler:      service.ServiceUpdateH,
			data:         sm5,
			err:          "{\"code\":400,\"status\":\"Bad Request\",\"message\":\"resources cpu limit exceeded\"}",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
		},
		// TODO: check another spec parameters
		{
			name:         "check update service success",
			fields:       fields{stg},
			args:         args{ctx, ns1, s1},
			handler:      service.ServiceUpdateH,
			data:         m3,
			want:         v1.View().Service().NewWithDeployment(s3, nil, nil),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "check update service success with limits",
			args:         args{ctx, ns3, s4},
			fields:       fields{stg},
			handler:      service.ServiceUpdateH,
			data:         sm6,
			want:         v1.View().Service().NewWithDeployment(s4, nil, nil),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "check update service success with equal limits",
			args:         args{ctx, ns3, s4},
			fields:       fields{stg},
			handler:      service.ServiceUpdateH,
			data:         sm7,
			want:         v1.View().Service().NewWithDeployment(s4, nil, nil),
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Service(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns3.SelfLink().String(), ns3, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Service(), s1.SelfLink().String(), s1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Service(), s4.SelfLink().String(), s4, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			bd, err := tc.data.ToJson()
			assert.NoError(t, err)

			req, err := http.NewRequest("PUT", fmt.Sprintf("/namespace/%s/service/%s", tc.args.namespace.Meta.Name, tc.args.service.Meta.Name), strings.NewReader(string(bd)))
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)
				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/namespace/{namespace}/service/{service}", tc.handler)

			setRequestVars(r, req)

			// We create assert ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			res := httptest.NewRecorder()

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			r.ServeHTTP(res, req)

			// Check the status code is what we expect.
			if !assert.Equal(t, tc.expectedCode, res.Code, "status code not equal") {
				return
			}

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			if tc.wantErr && res.Code != 200 {
				assert.Equal(t, tc.err, string(body), "incorrect status code")
			} else {
				s := new(views.Service)
				err := json.Unmarshal(body, &s)
				assert.NoError(t, err)

				assert.Equal(t, tc.want.Meta.Name, s.Meta.Name, "description not equal")
				assert.Equal(t, tc.want.Meta.Description, s.Meta.Description, "description not equal")
				assert.Equal(t, tc.want.Spec.Replicas, s.Spec.Replicas, "replicas not equal")
				assert.Equal(t, tc.want.Spec.Network.IP, s.Spec.Network.IP, "network ip spec not equal")
				assert.Equal(t, tc.want.Spec.Network.Ports, s.Spec.Network.Ports, "network ports spec not equal")
				assert.Equal(t, tc.want.Spec.Strategy.Type, s.Spec.Strategy.Type, "deployment strategy not equal")
				assert.Equal(t, tc.want.Spec.Selector.Node, s.Spec.Selector.Node, "provision node selectors not equal")
				assert.Equal(t, len(tc.want.Spec.Selector.Labels), len(s.Spec.Selector.Labels), "provision labels selectors not equal")

				assert.Equal(t, len(tc.want.Spec.Template.Containers), len(s.Spec.Template.Containers), "container spec count not equal")

				for _, wcs := range tc.want.Spec.Template.Containers {
					var f = false

					for _, scs := range s.Spec.Template.Containers {

						if scs.Name != wcs.Name {
							continue
						}

						f = true

						assert.Equal(t, wcs.Command, scs.Command, "container spec command not equal")
						assert.Equal(t, wcs.Entrypoint, scs.Entrypoint, "container spec entrypoint not equal")
						assert.Equal(t, wcs.Workdir, scs.Workdir, "container spec workdir not equal")

						assert.Equal(t, strings.Join(wcs.Args, " "), strings.Join(scs.Args, " "), "container spec command args not equal")

						assert.Equal(t, wcs.Resources, scs.Resources, "container resources not equal")
						assert.Equal(t, wcs.Image, scs.Image, "container spec image not equal")

						for _, wvcs := range wcs.Volumes {
							var vf = false

							for _, svcs := range scs.Volumes {

								if wvcs.Name != svcs.Name {
									continue
								}

								vf = true

								assert.Equal(t, wvcs, svcs, "container volume spec not equal")

							}

							if !vf {
								t.Error("container volume not found", wcs.Name)
							}

						}

						for _, wecs := range wcs.Env {

							var ef = false

							for _, secs := range scs.Env {
								if wecs.Name != secs.Name {
									continue
								}

								ef = true

								assert.Equal(t, wecs, secs, "container env spec not equal")

							}

							if !ef {
								t.Error("container env not found", wecs.Name)
								return
							}

						}

						assert.Equal(t, len(wcs.Env), len(scs.Env), "container count spec envs not equal")
					}

					if !f {
						assert.Error(t, errors.New("container spec not found"), wcs.Name)
						return
					}
				}

				if !assert.Equal(t, len(tc.want.Spec.Template.Volumes), len(s.Spec.Template.Volumes), "volumes specs count not equal") {
					return
				}

				for _, wvs := range tc.want.Spec.Template.Volumes {

					var f = false

					for _, scs := range s.Spec.Template.Volumes {

						if scs.Name != wvs.Name {
							continue
						}

						f = true

						assert.Equal(t, wvs.Type, scs.Type, "volume spec type not equal")
						assert.Equal(t, wvs.Secret.Name, scs.Secret.Name, "volume spec secret name not equal")

						assert.Equal(t, len(wvs.Secret.Binds), len(scs.Secret.Binds), "container spec secret binds not equal")

					}

					if !f {
						t.Log("not found")
						t.Error("volume spec not found", wvs.Name)
						return
					}

				}
			}
		})
	}

}

// Testing ServiceRemoveH handler
func TestServiceRemove(t *testing.T) {

	var ctx = context.Background()

	v := viper.New()
	v.SetDefault("storage.driver", "mock")

	stg, _ := storage.Get(v)
	envs.Get().SetStorage(stg)

	ns1 := getNamespaceAsset("demo", "")
	ns2 := getNamespaceAsset("test", "")

	s1 := getServiceAsset(ns1.Meta.Name, "demo", "")
	s2 := getServiceAsset(ns1.Meta.Name, "test", "")

	type fields struct {
		stg storage.Storage
	}

	type args struct {
		ctx       context.Context
		namespace *types.Namespace
		service   *types.Service
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		headers      map[string]string
		handler      func(http.ResponseWriter, *http.Request)
		err          string
		want         string
		wantErr      bool
		expectedCode int
	}{
		{
			name:         "checking get service if not exists",
			fields:       fields{stg},
			args:         args{ctx, ns1, s2},
			handler:      service.ServiceRemoveH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Service not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get service if namespace not exists",
			fields:       fields{stg},
			args:         args{ctx, ns2, s1},
			handler:      service.ServiceRemoveH,
			err:          "{\"code\":404,\"status\":\"Not Found\",\"message\":\"Namespace not found\"}",
			wantErr:      true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "checking get service successfully",
			fields:       fields{stg},
			args:         args{ctx, ns1, s1},
			handler:      service.ServiceRemoveH,
			want:         "",
			wantErr:      false,
			expectedCode: http.StatusOK,
		},
	}

	clear := func() {
		err := envs.Get().GetStorage().Del(context.Background(), stg.Collection().Namespace(), types.EmptyString)
		assert.NoError(t, err)

		err = envs.Get().GetStorage().Del(context.Background(), stg.Collection().Service(), types.EmptyString)
		assert.NoError(t, err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			clear()
			defer clear()

			err := tc.fields.stg.Put(context.Background(), stg.Collection().Namespace(), ns1.SelfLink().String(), ns1, nil)
			assert.NoError(t, err)

			err = tc.fields.stg.Put(context.Background(), stg.Collection().Service(), s1.SelfLink().String(), s1, nil)
			assert.NoError(t, err)

			// Create assert request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest("DELETE", fmt.Sprintf("/namespace/%s/service/%s", tc.args.namespace.Meta.Name, tc.args.service.Meta.Name), nil)
			assert.NoError(t, err)

			if tc.headers != nil {
				for key, val := range tc.headers {
					req.Header.Set(key, val)

				}
			}

			r := mux.NewRouter()
			r.HandleFunc("/namespace/{namespace}/service/{service}", tc.handler)

			setRequestVars(r, req)

			// We create assert ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			res := httptest.NewRecorder()

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			r.ServeHTTP(res, req)

			// Check the status code is what we expect.
			if !assert.Equal(t, tc.expectedCode, res.Code, "status code not equal") {
				return
			}

			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)

			if tc.wantErr && res.Code != 200 {
				assert.Equal(t, tc.err, string(body), "incorrect status code")
			} else {

				got := new(types.Service)
				err := tc.fields.stg.Get(tc.args.ctx, stg.Collection().Service(), tc.args.service.SelfLink().String(), got, nil)
				if err != nil && !errors.Storage().IsErrEntityNotFound(err) {
					assert.NoError(t, err)
				}

				if got != nil {
					assert.Equal(t, got.Status.State, types.StateDestroy, "status not destroy")
				}

				assert.Equal(t, tc.want, string(body), "response not equal with want")
			}

		})
	}

}

func getNamespaceAsset(name, desc string) *types.Namespace {
	var n = types.Namespace{}
	n.Meta.SetDefault()
	n.Meta.Name = name
	n.Meta.Description = desc
	n.Meta.SelfLink = *types.NewNamespaceSelfLink(name)
	return &n
}

func getServiceAsset(namespace, name, desc string) *types.Service {
	var s = types.Service{}
	s.Meta.SetDefault()
	s.Meta.Namespace = namespace
	s.Meta.Name = name
	s.Meta.Description = desc
	s.Spec.Replicas = 1
	s.Spec.Template.Containers = make(types.SpecTemplateContainers, 0)
	s.Spec.Template.Containers = append(s.Spec.Template.Containers, &types.SpecTemplateContainer{
		Name: "demo",
	})
	s.Meta.SelfLink = *types.NewServiceSelfLink(namespace, name)
	return &s
}

func setRequestVars(r *mux.Router, req *http.Request) {
	var match mux.RouteMatch
	// Take the request and match it
	r.Match(req, &match)
	// Push the variable onto the context
	req = mux.SetURLVars(req, match.Vars)
}

func getServiceManifest(name, image string) *request.ServiceManifest {

	var (
		replicas  = 1
		container = request.ManifestSpecTemplateContainer{
			Name: image,
			Image: &request.ManifestSpecTemplateContainerImage{
				Name: image,
			},
			Env: make([]request.ManifestSpecTemplateContainerEnv, 0),
		}
		volume = request.ManifestSpecTemplateVolume{
			Name: "demo",
			Secret: &request.ManifestSpecTemplateSecretVolume{
				Name:  "test",
				Binds: make([]request.ManifestSpecTemplateSecretVolumeBind, 0),
			},
		}
	)

	volume.Secret.Binds = append(volume.Secret.Binds, request.ManifestSpecTemplateSecretVolumeBind{
		Key:  "demo",
		File: "test.txt",
	})

	container.Env = append(container.Env, request.ManifestSpecTemplateContainerEnv{
		Name:  "Demo",
		Value: "test",
	})

	container.Env = append(container.Env, request.ManifestSpecTemplateContainerEnv{
		Name: "Secret",
		Secret: &request.ManifestSpecTemplateContainerEnvSecret{
			Name: "secret-name",
			Key:  "secret-key",
		},
	})

	mf := new(request.ServiceManifest)
	mf.Meta.Name = &name
	mf.Spec.Replicas = &replicas
	mf.Spec.Template = new(request.ManifestSpecTemplate)
	mf.Spec.Template.Containers = append(mf.Spec.Template.Containers, container)
	mf.Spec.Template.Volumes = append(mf.Spec.Template.Volumes, volume)
	return mf
}
