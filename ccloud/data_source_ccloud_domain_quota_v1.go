package ccloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceCCloudDomainQuotaV1() *schema.Resource {
	quotaResource := &schema.Resource{
		Read: dataSourceCCloudDomainQuotaV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}

	for service, resources := range limesServices {
		elem := &schema.Resource{
			Schema: make(map[string]*schema.Schema, len(resources)),
		}

		for resource := range resources {
			elem.Schema[resource] = &schema.Schema{
				Type:     schema.TypeFloat,
				Computed: true,
			}
		}

		quotaResource.Schema[sanitize(service)] = &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem:     elem,
			MaxItems: 1,
		}
	}

	return quotaResource
}

func dataSourceCCloudDomainQuotaV1Read(d *schema.ResourceData, meta interface{}) error {
	domainID := d.Get("domain_id").(string)
	d.SetId(domainID)

	return resourceCCloudDomainQuotaV1Read(d, meta)
}
