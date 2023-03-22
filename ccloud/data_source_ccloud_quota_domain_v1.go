package ccloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCCloudQuotaDomainV1Deprecated() *schema.Resource {
	quotaResource := &schema.Resource{
		DeprecationMessage: "use ccloud_quota_domain_v1 data source instead",
		ReadContext:        dataSourceCCloudQuotaDomainV1Read,

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
		}
	}

	return quotaResource
}

func dataSourceCCloudQuotaDomainV1() *schema.Resource {
	quotaResource := &schema.Resource{
		ReadContext: dataSourceCCloudQuotaDomainV1Read,

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
		}
	}

	return quotaResource
}

func dataSourceCCloudQuotaDomainV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	domainID := d.Get("domain_id").(string)
	d.SetId(domainID)

	return resourceCCloudQuotaDomainV1Read(ctx, d, meta)
}
