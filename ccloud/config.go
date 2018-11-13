package ccloud

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/hashicorp/terraform/helper/pathorcontents"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sapcc/gophercloud-limes/resources"
)

type Config struct {
	CACertFile        string
	ClientCertFile    string
	ClientKeyFile     string
	Cloud             string
	UserDomainName    string
	UserDomainID      string
	ProjectDomainName string
	ProjectDomainID   string
	DomainID          string
	DomainName        string
	EndpointType      string
	IdentityEndpoint  string
	Insecure          *bool
	Password          string
	Region            string
	Swauth            bool
	TenantID          string
	TenantName        string
	Token             string
	Username          string
	UserID            string
	useOctavia        bool

	OsClient *gophercloud.ProviderClient
}

func (c *Config) LoadAndValidate() error {
	log.Printf("[CCLOUD] LoadAndValidate")
	c.Debug()

	// Make sure at least one of auth_url or cloud was specified.
	if c.IdentityEndpoint == "" && c.Cloud == "" {
		return fmt.Errorf("One of 'auth_url' or 'cloud' must be specified")
	}

	validEndpoint := false
	validEndpoints := []string{
		"internal", "internalURL",
		"admin", "adminURL",
		"public", "publicURL",
		"",
	}

	for _, endpoint := range validEndpoints {
		if c.EndpointType == endpoint {
			validEndpoint = true
		}
	}

	if !validEndpoint {
		return fmt.Errorf("Invalid endpoint type provided")
	}

	clientOpts := new(clientconfig.ClientOpts)

	// If a cloud entry was given, base AuthOptions on a clouds.yaml file.
	if c.Cloud != "" {
		clientOpts.Cloud = c.Cloud

		cloud, err := clientconfig.GetCloudFromYAML(clientOpts)
		if err != nil {
			return err
		}

		if c.Region == "" && cloud.RegionName != "" {
			c.Region = cloud.RegionName
		}

		if c.CACertFile == "" && cloud.CACertFile != "" {
			c.CACertFile = cloud.CACertFile
		}

		if c.ClientCertFile == "" && cloud.ClientCertFile != "" {
			c.ClientCertFile = cloud.ClientCertFile
		}

		if c.ClientKeyFile == "" && cloud.ClientKeyFile != "" {
			c.ClientKeyFile = cloud.ClientKeyFile
		}

		if c.Insecure == nil && cloud.Verify != nil {
			v := (!*cloud.Verify)
			c.Insecure = &v
		}
	} else {
		authInfo := &clientconfig.AuthInfo{
			AuthURL:           c.IdentityEndpoint,
			DomainID:          c.DomainID,
			DomainName:        c.DomainName,
			Password:          c.Password,
			ProjectDomainID:   c.ProjectDomainID,
			ProjectDomainName: c.ProjectDomainName,
			ProjectID:         c.TenantID,
			ProjectName:       c.TenantName,
			Token:             c.Token,
			UserDomainID:      c.UserDomainID,
			UserDomainName:    c.UserDomainName,
			Username:          c.Username,
			UserID:            c.UserID,
		}
		clientOpts.AuthInfo = authInfo
	}

	ao, err := clientconfig.AuthOptions(clientOpts)
	if err != nil {
		return err
	}

	client, err := openstack.NewClient(ao.IdentityEndpoint)
	if err != nil {
		return err
	}

	// Set UserAgent
	client.UserAgent.Prepend(terraform.UserAgentString())

	config := &tls.Config{}
	if c.CACertFile != "" {
		caCert, _, err := pathorcontents.Read(c.CACertFile)
		if err != nil {
			return fmt.Errorf("Error reading CA Cert: %s", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(caCert))
		config.RootCAs = caCertPool
	}

	if c.Insecure == nil {
		config.InsecureSkipVerify = false
	} else {
		config.InsecureSkipVerify = *c.Insecure
	}

	if c.ClientCertFile != "" && c.ClientKeyFile != "" {
		clientCert, _, err := pathorcontents.Read(c.ClientCertFile)
		if err != nil {
			return fmt.Errorf("Error reading Client Cert: %s", err)
		}
		clientKey, _, err := pathorcontents.Read(c.ClientKeyFile)
		if err != nil {
			return fmt.Errorf("Error reading Client Key: %s", err)
		}

		cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
		if err != nil {
			return err
		}

		config.Certificates = []tls.Certificate{cert}
		config.BuildNameToCertificate()
	}

	transport := &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: config}
	client.HTTPClient = http.Client{
		Transport: transport,
	}

	// If using Swift Authentication, there's no need to validate authentication normally.
	if !c.Swauth {
		err = openstack.Authenticate(client, *ao)
		if err != nil {
			return err
		}
	}

	c.OsClient = client

	return nil
}

func (c *Config) determineRegion(region string) string {
	// If a resource-level region was not specified, and a provider-level region was set,
	// use the provider-level region.
	if region == "" && c.Region != "" {
		region = c.Region
	}

	log.Printf("[DEBUG] OpenStack Region is: %s", region)
	return region
}

func (c *Config) getEndpointType() gophercloud.Availability {
	if c.EndpointType == "internal" || c.EndpointType == "internalURL" {
		return gophercloud.AvailabilityInternal
	}
	if c.EndpointType == "admin" || c.EndpointType == "adminURL" {
		return gophercloud.AvailabilityAdmin
	}
	return gophercloud.AvailabilityPublic
}

func (c *Config) limesV1Client(region string) (*gophercloud.ServiceClient, error) {
	c.Debug()

	return resources.NewLimesV1(c.OsClient, gophercloud.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getEndpointType(),
	})
}

func (c *Config) kubernikusV1Client(region string, isAdmin bool) (*Kubernikus, error) {
	c.Debug()

	serviceType := "kubernikus"
	if isAdmin {
		serviceType = "kubernikus-kubernikus"
	}

	return NewKubernikusV1(c.OsClient, gophercloud.EndpointOpts{
		Type:         serviceType,
		Region:       c.determineRegion(region),
		Availability: gophercloud.AvailabilityPublic,
	})
}

func (c *Config) identityV3Client(region string) (*gophercloud.ServiceClient, error) {
	return openstack.NewIdentityV3(c.OsClient, gophercloud.EndpointOpts{
		Region:       c.determineRegion(region),
		Availability: c.getEndpointType(),
	})
}

func (c *Config) Debug() {
	log.Printf("[CCLOUD] cacert_file:         %s", c.CACertFile)
	log.Printf("[CCLOUD] cert:                %s", c.ClientCertFile)
	log.Printf("[CCLOUD] key:                 %s", c.ClientKeyFile)
	log.Printf("[CCLOUD] cloud:               %s", c.Cloud)
	log.Printf("[CCLOUD] domain_id:           %s", c.DomainID)
	log.Printf("[CCLOUD] domain_name:         %s", c.DomainName)
	log.Printf("[CCLOUD] endpoint_type:       %s", c.EndpointType)
	log.Printf("[CCLOUD] auth_url:            %s", c.IdentityEndpoint)
	log.Printf("[CCLOUD] password:            %s", c.Password)
	log.Printf("[CCLOUD] project_domain_id:   %s", c.ProjectDomainID)
	log.Printf("[CCLOUD] project_domain_name: %s", c.ProjectDomainName)
	log.Printf("[CCLOUD] region:              %s", c.Region)
	log.Printf("[CCLOUD] swauth:              %s", c.Swauth)
	log.Printf("[CCLOUD] token:               %s", c.Token)
	log.Printf("[CCLOUD] tenant_id:           %s", c.TenantID)
	log.Printf("[CCLOUD] tenant_name:         %s", c.TenantName)
	log.Printf("[CCLOUD] user_domain_id:      %s", c.UserDomainID)
	log.Printf("[CCLOUD] user_domain_name:    %s", c.UserDomainName)
	log.Printf("[CCLOUD] user_name:           %s", c.Username)
	log.Printf("[CCLOUD] user_id:             %s", c.UserID)
	log.Printf("[CCLOUD] use_octavia:         %s", c.useOctavia)
}
