package ccloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/sapcc/gophercloud-sapcc/arc/v1/agents"

	"github.com/gophercloud/utils/terraform/hashcode"
)

func resourceCCloudArcAgentBootstrapV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCCloudArcAgentBootstrapV1Create,
		ReadContext:   func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics { return nil },
		DeleteContext: func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics { return nil },

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "cloud-config",
				ValidateFunc: validation.StringInSlice([]string{
					"linux", "windows", "cloud-config", "json",
				}, false),
			},

			"triggers": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			// computed attributes
			"user_data": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"raw_map": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func resourceCCloudArcAgentBootstrapV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	arcClient, err := config.arcV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack Arc client: %s", err)
	}

	var bootstrapType string
	switch d.Get("type") {
	case "linux":
		bootstrapType = "text/x-shellscript"
	case "windows":
		bootstrapType = "text/x-powershellscript"
	case "cloud-config":
		bootstrapType = "text/cloud-config"
	case "json":
		bootstrapType = "application/json"
	}

	createOpts := agents.InitOpts{Accept: bootstrapType}

	log.Printf("[DEBUG] ccloud_arc_agent_bootstrap_v1 create options: %#v", createOpts)

	res := agents.Init(arcClient, createOpts)
	if res.Err != nil {
		return diag.Errorf("Error creating ccloud_arc_agent_bootstrap_v1: %s", res.Err)
	}

	headers, err := res.ExtractHeaders()
	if err != nil {
		return diag.Errorf("Error extracting headers while creating ccloud_arc_agent_bootstrap_v1: %s", err)
	}

	if bootstrapType != headers.ContentType {
		return diag.Errorf("Error verifying headers while creating ccloud_arc_agent_bootstrap_v1: wants '%s', got '%s'", bootstrapType, headers.ContentType)
	}

	data, err := res.ExtractContent()
	if err != nil {
		return diag.Errorf("Error extracting content while creating ccloud_arc_agent_bootstrap_v1: %s", err)
	}

	userData := string(data)

	d.SetId(fmt.Sprintf("%d", hashcode.String(userData)))

	if bootstrapType == "application/json" {
		var initMap map[string]string
		err = json.Unmarshal(data, &initMap)
		if err != nil {
			return diag.Errorf("Error unmarshalling JSON content while creating ccloud_arc_agent_bootstrap_v1: %s", err)
		}
		d.Set("raw_map", initMap)
	} else {
		d.Set("raw_map", map[string]string{})
	}

	d.Set("user_data", userData)
	d.Set("region", GetRegion(d, config))

	return nil
}
