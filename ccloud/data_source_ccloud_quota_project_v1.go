package ccloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCCloudQuotaProjectV1Deprecated() *schema.Resource {
	quotaResource := &schema.Resource{
		DeprecationMessage: "use ccloud_quota_project_v1 data source instead",
		ReadContext:        dataSourceCCloudQuotaProjectV1Read,

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
			"bursting": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"multiplier": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
					},
				},
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

func dataSourceCCloudQuotaProjectV1() *schema.Resource {
	quotaResource := &schema.Resource{
		ReadContext: dataSourceCCloudQuotaProjectV1Read,

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
			"bursting": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"multiplier": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
					},
				},
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

func dataSourceCCloudQuotaProjectV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectID := d.Get("project_id").(string)

	d.SetId(projectID)

	return resourceCCloudQuotaProjectV1Read(ctx, d, meta)
}
