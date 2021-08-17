package ccloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceCCloudProjectQuotaV1() *schema.Resource {
	quotaResource := &schema.Resource{
		Read: dataSourceCCloudProjectQuotaV1Read,

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
			"project_id": {
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

func dataSourceCCloudProjectQuotaV1Read(d *schema.ResourceData, meta interface{}) error {
	projectID := d.Get("project_id").(string)

	d.SetId(projectID)

	return resourceCCloudProjectQuotaV1Read(d, meta)
}
