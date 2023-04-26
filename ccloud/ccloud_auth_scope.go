package ccloud

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
)

type authScopeTokenDetails struct {
	user    *tokens.User
	domain  *tokens.Domain
	project *tokens.Project
	catalog *tokens.ServiceCatalog
	roles   []tokens.Role
}

func getTokenDetails(sc *gophercloud.ServiceClient) (authScopeTokenDetails, error) {
	var (
		details authScopeTokenDetails
		err     error
	)

	r := sc.ProviderClient.GetAuthResult()
	switch result := r.(type) {
	case tokens.CreateResult:
		details.user, err = result.ExtractUser()
		if err != nil {
			return details, err
		}
		details.domain, err = result.ExtractDomain()
		if err != nil {
			return details, err
		}
		details.project, err = result.ExtractProject()
		if err != nil {
			return details, err
		}
		details.roles, err = result.ExtractRoles()
		if err != nil {
			return details, err
		}
		details.catalog, err = result.ExtractServiceCatalog()
		if err != nil {
			return details, err
		}
	case tokens.GetResult:
		details.user, err = result.ExtractUser()
		if err != nil {
			return details, err
		}
		details.domain, err = result.ExtractDomain()
		if err != nil {
			return details, err
		}
		details.project, err = result.ExtractProject()
		if err != nil {
			return details, err
		}
		details.roles, err = result.ExtractRoles()
		if err != nil {
			return details, err
		}
		details.catalog, err = result.ExtractServiceCatalog()
		if err != nil {
			return details, err
		}
	default:
		res := tokens.Get(sc, sc.ProviderClient.TokenID)
		if res.Err != nil {
			return details, res.Err
		}
		details.user, err = res.ExtractUser()
		if err != nil {
			return details, err
		}
		details.domain, err = res.ExtractDomain()
		if err != nil {
			return details, err
		}
		details.project, err = res.ExtractProject()
		if err != nil {
			return details, err
		}
		details.roles, err = res.ExtractRoles()
		if err != nil {
			return details, err
		}
		details.catalog, err = res.ExtractServiceCatalog()
		if err != nil {
			return details, err
		}
	}

	return details, nil
}
