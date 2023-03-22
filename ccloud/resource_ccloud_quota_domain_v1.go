package ccloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	limesresources "github.com/sapcc/go-api-declarations/limes/resources"
	"github.com/sapcc/gophercloud-sapcc/resources/v1/domains"
)

func resourceCCloudQuotaDomainV1Deprecated() *schema.Resource {
	quotaResource := &schema.Resource{
		DeprecationMessage: "use ccloud_quota_domain_v1 resource instead",
		SchemaVersion:      1,

		ReadContext:   resourceCCloudQuotaDomainV1Read,
		UpdateContext: resourceCCloudQuotaDomainV1CreateOrUpdate,
		CreateContext: resourceCCloudQuotaDomainV1CreateOrUpdate,
		DeleteContext: resourceCCloudQuotaDomainV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCCloudQuotaDomainV1Import,
		},

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
				Required: false,
				Optional: true,
				Computed: true,
			}
		}

		quotaResource.Schema[sanitize(service)] = &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem:     elem,
			MaxItems: 1,
		}
	}

	return quotaResource
}

func resourceCCloudQuotaDomainV1() *schema.Resource {
	quotaResource := &schema.Resource{
		SchemaVersion: 1,

		ReadContext:   resourceCCloudQuotaDomainV1Read,
		UpdateContext: resourceCCloudQuotaDomainV1CreateOrUpdate,
		CreateContext: resourceCCloudQuotaDomainV1CreateOrUpdate,
		DeleteContext: resourceCCloudQuotaDomainV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCCloudQuotaDomainV1Import,
		},

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
				Required: false,
				Optional: true,
				Computed: true,
			}
		}

		quotaResource.Schema[sanitize(service)] = &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem:     elem,
			MaxItems: 1,
		}
	}

	return quotaResource
}

func resourceCCloudQuotaDomainV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	domainID := d.Get("domain_id").(string)

	log.Printf("[DEBUG] Reading Quota for: %s", domainID)

	config := meta.(*Config)
	limes, err := config.limesV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack limes client: %s", err)
	}

	quota, err := domains.Get(limes, domainID, domains.GetOpts{}).Extract()
	if err != nil {
		return diag.Errorf("Error getting Limes domain: %s", err)
	}

	for service, resources := range limesServices {
		res := make(map[string]*uint64)
		for resource := range resources {
			if quota.Services[service] == nil || quota.Services[service].Resources[resource] == nil {
				continue
			}
			res[resource] = quota.Services[service].Resources[resource].DomainQuota
			log.Printf("[DEBUG] %s.%s: %s", service, resource, toString(quota.Services[service].Resources[resource]))
		}
		d.Set(sanitize(service), []map[string]*uint64{res})
	}

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceCCloudQuotaDomainV1CreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	domainID := d.Get("domain_id").(string)
	services := limesresources.QuotaRequest{}

	log.Printf("[DEBUG] Updating Quota for: %s", domainID)

	config := meta.(*Config)
	client, err := config.limesV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack limes client: %s", err)
	}

	for _service, resources := range limesServices {
		service := sanitize(_service)
		if _, ok := d.GetOk(service); ok && d.HasChange(service) {
			log.Printf("[DEBUG] Service Changed: %s", service)

			quota := make(limesresources.ServiceQuotaRequest)
			for resource, unit := range resources {
				key := fmt.Sprintf("%s.0.%s", service, resource)

				if d.HasChange(key) {
					v := d.Get(key)
					log.Printf("[DEBUG] Resource Changed: %s", key)
					quota[resource] = limesresources.ResourceQuotaRequest{Value: uint64(v.(float64)), Unit: unit}
					log.Printf("[DEBUG] %s.%s: %v", service, resource, quota[resource])
				}
			}
			services[_service] = quota
		}
	}

	opts := domains.UpdateOpts{Services: services}
	err = domains.Update(client, domainID, opts).ExtractErr()
	if err != nil {
		return diag.Errorf("Error updating Limes domain: %s", err)
	}

	log.Printf("[DEBUG] Resulting Quota for: %s", domainID)

	d.SetId(domainID)

	return resourceCCloudQuotaDomainV1Read(ctx, d, meta)
}

func resourceCCloudQuotaDomainV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

func resourceCCloudQuotaDomainV1Import(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.Set("domain_id", d.Id())

	return []*schema.ResourceData{d}, nil
}
