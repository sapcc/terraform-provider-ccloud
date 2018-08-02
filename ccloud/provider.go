package ccloud

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var descriptions = map[string]string{
	"auth_url":  "The Identity authentication URL.",
	"region":    "The OpenStack region to connect to.",
	"user_name": "Username to login with.",
	"user_id":   "User ID to login with.",
	"tenant_id": "The ID of the Tenant (Identity v2) or Project (Identity v3)\n" +
		"to login with.",
	"tenant_name": "The name of the Tenant (Identity v2) or Project (Identity v3)\n" +
		"to login with.",
	"password":            "Password to login with.",
	"token":               "Authentication token to use as an alternative to username/password.",
	"user_domain_name":    "The name of the domain where the user resides (Identity v3).",
	"user_domain_id":      "The ID of the domain where the user resides (Identity v3).",
	"project_domain_name": "The name of the domain where the project resides (Identity v3).",
	"project_domain_id":   "The ID of the domain where the proejct resides (Identity v3).",
	"domain_id":           "The ID of the Domain to scope to (Identity v3).",
	"domain_name":         "The name of the Domain to scope to (Identity v3).",
	"insecure":            "Trust self-signed certificates.",
	"cacert_file":         "A Custom CA certificate.",
	"endpoint_type":       "The catalog endpoint type to use.",
	"cert":                "A client certificate to authenticate with.",
	"key":                 "A client private key to authenticate with.",
	"swauth": "Use Swift's authentication system instead of Keystone. Only used for\n" +
		"interaction with Swift.",
	"use_octavia": "If set to `true`, API requests will go the Load Balancer\n" +
		"service (Octavia) instead of the Networking service (Neutron).",
	"cloud": "An entry in a `clouds.yaml` file to use.",
}

func Provider() terraform.ResourceProvider {
	log.Printf("[CCLOUD] CCloud Provider Init")

	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_AUTH_URL", ""),
				Description: descriptions["auth_url"],
			},

			"region": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["region"],
				DefaultFunc: schema.EnvDefaultFunc("OS_REGION_NAME", ""),
			},

			"user_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USERNAME", ""),
				Description: descriptions["user_name"],
			},

			"user_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USER_ID", ""),
				Description: descriptions["user_name"],
			},

			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TENANT_ID",
					"OS_PROJECT_ID",
				}, ""),
				Description: descriptions["tenant_id"],
			},

			"tenant_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TENANT_NAME",
					"OS_PROJECT_NAME",
				}, ""),
				Description: descriptions["tenant_name"],
			},

			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OS_PASSWORD", ""),
				Description: descriptions["password"],
			},

			"token": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_TOKEN",
					"OS_AUTH_TOKEN",
				}, ""),
				Description: descriptions["token"],
			},

			"user_domain_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_USER_DOMAIN_NAME",
				}, ""),
				Description: descriptions["user_domain_name"],
			},

			"user_domain_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_USER_DOMAIN_ID",
				}, ""),
				Description: descriptions["user_domain_id"],
			},

			"project_domain_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_PROJECT_DOMAIN_NAME",
				}, ""),
				Description: descriptions["project_domain_name"],
			},

			"project_domain_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_PROJECT_DOMAIN_ID",
				}, ""),
				Description: descriptions["project_domain_id"],
			},

			"domain_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_DOMAIN_ID",
				}, ""),
				Description: descriptions["domain_id"],
			},

			"domain_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"OS_DOMAIN_NAME",
					"OS_DEFAULT_DOMAIN",
				}, ""),
				Description: descriptions["domain_name"],
			},

			"insecure": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_INSECURE", nil),
				Description: descriptions["insecure"],
			},

			"endpoint_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_ENDPOINT_TYPE", ""),
			},

			"cacert_file": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CACERT", ""),
				Description: descriptions["cacert_file"],
			},

			"cert": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CERT", ""),
				Description: descriptions["cert"],
			},

			"key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_KEY", ""),
				Description: descriptions["key"],
			},

			"swauth": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_SWAUTH", ""),
				Description: descriptions["swauth"],
			},

			"use_octavia": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_USE_OCTAVIA", ""),
				Description: descriptions["use_octavia"],
			},

			"cloud": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_CLOUD", ""),
				Description: descriptions["cloud"],
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"ccloud_quota":      resourceCCloudQuota(),
			"ccloud_kubernetes": resourceCCloudKubernetes(),
		},

		ConfigureFunc: configureProvider,
	}

}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		CACertFile:        d.Get("cacert_file").(string),
		ClientCertFile:    d.Get("cert").(string),
		ClientKeyFile:     d.Get("key").(string),
		Cloud:             d.Get("cloud").(string),
		DomainID:          d.Get("domain_id").(string),
		DomainName:        d.Get("domain_name").(string),
		EndpointType:      d.Get("endpoint_type").(string),
		IdentityEndpoint:  d.Get("auth_url").(string),
		Password:          d.Get("password").(string),
		ProjectDomainID:   d.Get("project_domain_id").(string),
		ProjectDomainName: d.Get("project_domain_name").(string),
		Region:            d.Get("region").(string),
		Swauth:            d.Get("swauth").(bool),
		Token:             d.Get("token").(string),
		TenantID:          d.Get("tenant_id").(string),
		TenantName:        d.Get("tenant_name").(string),
		UserDomainID:      d.Get("user_domain_id").(string),
		UserDomainName:    d.Get("user_domain_name").(string),
		Username:          d.Get("user_name").(string),
		UserID:            d.Get("user_id").(string),
		useOctavia:        d.Get("use_octavia").(bool),
	}

	log.Printf("[CCLOUD] %v", config)

	v, ok := d.GetOk("insecure")
	if ok {
		insecure := v.(bool)
		config.Insecure = &insecure
	}

	if err := config.LoadAndValidate(); err != nil {
		return nil, err
	}

	return &config, nil
}
