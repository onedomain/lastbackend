//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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

package routes

import (
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	"github.com/lastbackend/lastbackend/pkg/daemon/config"
	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/errors"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"github.com/lastbackend/lastbackend/pkg/vendors"
	"github.com/lastbackend/lastbackend/pkg/vendors/interfaces"
	"net/http"
)

// Авторизация сторонних сервисов для платформы
func OAuthConnectH(w http.ResponseWriter, r *http.Request) {

	var (
		log            = c.Get().GetLogger()
		storage        = c.Get().GetStorage()
		clientID       string
		clientSecretID string
		redirectURI    string
		client         interfaces.IAuth
		params         = utils.Vars(r)
		vendor         = params[`vendor`]
		code           = params[`code`]
	)

	log.Debug("Connect service handler")

	clientID, clientSecretID, redirectURI = config.Get().GetVendorConfig(vendor)

	if clientID == "" || clientSecretID == "" {
		log.Error("Error: user unauthorized")
		errors.HTTP.Unauthorized(w)
		return
	}

	// Get client for github/bitbucket/gitlab (or anything if implement adapter.OAuth interface) by vendor in user or organization mode
	switch vendor {
	case "github":
		client = vendors.GetGitHub(clientID, clientSecretID, redirectURI)
	case "bitbucket":
		client = vendors.GetBitBucket(clientID, clientSecretID, redirectURI)
	case "gitlab":
		client = vendors.GetGitLab(clientID, clientSecretID, redirectURI)
	default:
		log.Error("vendor is not supported yet")
		errors.BadParameter("vendor").Http(w)
		return
	}

	token, err := client.GetToken(code)
	if err != nil {
		log.Error(err)
		errors.HTTP.InternalServerError(w)
		return
	}

	serviceUser, err := client.GetUser(token)
	if err != nil {
		log.Error(err)
		errors.HTTP.InternalServerError(w)
		return
	}

	vendorInfo := client.GetVendorInfo()

	if err := storage.Vendor().Insert(r.Context(), serviceUser.Username, vendorInfo.Vendor, vendorInfo.Host, serviceUser.ServiceID, token); err != nil {
		log.Error(err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})
}

func OAuthDisconnectH(w http.ResponseWriter, r *http.Request) {

	var (
		log            = c.Get().GetLogger()
		storage        = c.Get().GetStorage()
		params = utils.Vars(r)
		vendor = params[`vendor`]
	)

	log.Debug("Disconnect service handler")

	if err := storage.Vendor().Remove(r.Context(), vendor); err != nil {
		log.Error(err)
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})
}

func VCSRepositoriesListH(w http.ResponseWriter, r *http.Request) {

	var (
		log            = c.Get().GetLogger()
		storage        = c.Get().GetStorage()
		client interfaces.IVCS
		params = utils.Vars(r)
		vendor = params[`vendor`]
	)

	clientID, clientSecretID, redirectURI := config.Get().GetVendorConfig(vendor)

	if clientID == "" || clientSecretID == "" {
		log.Error("Error: user unauthorized")
		errors.HTTP.Unauthorized(w)
		return
	}

	// Get client for github/bitbucket/gitlab (or anything if implement adapter.OAuth interface) by vendor in user or organization mode
	switch vendor {
	case "github":
		client = vendors.GetGitHub(clientID, clientSecretID, redirectURI)
	case "bitbucket":
		client = vendors.GetBitBucket(clientID, clientSecretID, redirectURI)
	case "gitlab":
		client = vendors.GetGitLab(clientID, clientSecretID, redirectURI)
	default:
		log.Error("vendor is not supported yet")
		errors.BadParameter("vendor").Http(w)
		return
	}

	// ************************ Update token ************************ //

	vendorInfo := client.GetVendorInfo()

	oaModel, err := storage.Vendor().Get(r.Context(), vendorInfo.Vendor)
	if err != nil {
		log.Error(err)
		errors.HTTP.InternalServerError(w)
	}

	token, modify, err := client.RefreshToken(oaModel.Token)
	if err != nil {
		log.Error(err)
		errors.HTTP.InternalServerError(w)
	}

	u, err := client.GetUser(token)
	if err != nil {
		log.Error(err)
		errors.HTTP.InternalServerError(w)
	}

	if modify {

		oaModel.Host = vendorInfo.Host
		oaModel.Vendor = vendorInfo.Vendor
		oaModel.ServiceID = u.ServiceID
		oaModel.Token = token
		oaModel.Username = u.Username

		if err = storage.Vendor().Update(r.Context(), oaModel); err != nil {
			log.Error(err)
			errors.HTTP.InternalServerError(w)
		}
	}

	// ************************ End update token ************************ //

	repos, err := client.ListRepositories(token, u.Username, false)
	if err != nil {
		log.Error(err)
		errors.HTTP.InternalServerError(w)
		return
	}

	rp := types.VCSRepositoryList{}
	rp.Convert(repos)
	response, err := rp.ToJson()

	w.WriteHeader(200)
	w.Write(response)
}

func VCSBranchesListH(w http.ResponseWriter, r *http.Request) {

	var (
		log            = c.Get().GetLogger()
		storage        = c.Get().GetStorage()
		client interfaces.IVCS
		params = utils.Vars(r)
		vendor = params[`vendor`]
		repo   = r.URL.Query().Get(`repo`)
	)

	clientID, clientSecretID, redirectURI := config.Get().GetVendorConfig(vendor)

	if clientID == "" || clientSecretID == "" {
		log.Error("Error: user unauthorized")
		errors.HTTP.Unauthorized(w)
		return
	}

	// Get client for github/bitbucket/gitlab (or anything if implement adapter.OAuth interface) by vendor in user or organization mode
	switch vendor {
	case "github":
		client = vendors.GetGitHub(clientID, clientSecretID, redirectURI)
	case "bitbucket":
		client = vendors.GetBitBucket(clientID, clientSecretID, redirectURI)
	case "gitlab":
		client = vendors.GetGitLab(clientID, clientSecretID, redirectURI)
	default:
		log.Error("vendor is not supported yet")
		errors.BadParameter("vendor").Http(w)
		return
	}

	// ************************ Update token ************************ //

	vendorInfo := client.GetVendorInfo()

	oaModel, err := storage.Vendor().Get(r.Context(), vendorInfo.Vendor)
	if err != nil {
		log.Error(err)
		errors.HTTP.InternalServerError(w)
	}

	token, modify, err := client.RefreshToken(oaModel.Token)
	if err != nil {
		log.Error(err)
		errors.HTTP.InternalServerError(w)
	}

	u, err := client.GetUser(token)
	if err != nil {
		log.Error(err)
		errors.HTTP.InternalServerError(w)
	}

	if modify {
		oaModel.Host = vendorInfo.Host
		oaModel.Vendor = vendorInfo.Vendor
		oaModel.ServiceID = u.ServiceID
		oaModel.Token = token
		oaModel.Username = u.Username

		if err = storage.Vendor().Update(r.Context(), oaModel); err != nil {
			log.Error(err)
			errors.HTTP.InternalServerError(w)
		}
	}

	// ************************ End update token ************************ //

	branches, err := client.ListBranches(token, u.Username, repo)
	if err != nil {
		log.Error(err)
		errors.HTTP.InternalServerError(w)
		return
	}

	br := types.VCSBranchList{}
	br.Convert(branches)
	response, err := br.ToJson()

	w.WriteHeader(200)
	w.Write(response)
}
