package ccloud

import (
	"log"
	"net/url"
	"reflect"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/gophercloud/gophercloud"
	"github.com/pkg/errors"
	"github.com/sapcc/kubernikus/pkg/api/client/operations"

	httptransport "github.com/go-openapi/runtime/client"
)

type Kubernikus struct {
	operations.Client
	provider *gophercloud.ProviderClient
}

func NewKubernikusV1(provider *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*Kubernikus, error) {
	var err error
	var endpoint string
	var kurl *url.URL

	if !reflect.DeepEqual(eo, gophercloud.EndpointOpts{}) {
		eo.ApplyDefaults("kubernikus")
		endpoint, err = provider.EndpointLocator(eo)
		if err != nil {
			return nil, err
		}
	}

	if kurl, err = url.Parse(endpoint); err != nil {
		return nil, errors.Errorf("Parsing the Kubernikus URL failed: %s", err)
	}

	transport := httptransport.New(kurl.Host, kurl.EscapedPath(), []string{kurl.Scheme})
	operations := operations.New(transport, strfmt.Default)

	return &Kubernikus{*operations, provider}, nil
}

func (k *Kubernikus) authFunc() runtime.ClientAuthInfoWriterFunc {
	return runtime.ClientAuthInfoWriterFunc(
		func(req runtime.ClientRequest, reg strfmt.Registry) error {
			log.Printf("[KUBERNETES] Kubernikus Auth %s", k.provider.Token())
			req.SetHeaderParam("X-AUTH-TOKEN", k.provider.Token())
			return nil
		})
}
