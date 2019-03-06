package ccloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/kayrus/gophercloud-arc/arc/v1/agents"
)

func resourceCCloudArcAgentV1() *schema.Resource {
	return &schema.Resource{
		Read:   resourceCCloudArcAgentV1Read,
		Create: resourceCCloudArcAgentV1Create,
		Update: resourceCCloudArcAgentV1Update,
		Delete: resourceCCloudArcAgentV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"filter": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			// Computed attributes
			"agent_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"project": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"organization": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated_with": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"all_tags": {
				Type:     schema.TypeMap,
				Computed: true,
			},

			"facts": {
				Type:     schema.TypeMap,
				Computed: true,
			},

			"facts_agents": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func resourceCCloudArcAgentV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	arcClient, err := config.arcV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack Arc client: %s", err)
	}

	var tmp interface{}
	filter := d.Get("filter").(string)
	timeout := d.Timeout(schema.TimeoutCreate)

	waitForAgent := &resource.StateChangeConf{
		Target:     []string{"active"},
		Refresh:    arcCCloudArcAgentV1GetAgent(arcClient, "", filter),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	tmp, err = waitForAgent.WaitForState()
	if err != nil {
		return err
	}

	agent := tmp.(*agents.Agent)

	d.SetId(agent.AgentID)

	err = updateArcAgentTagsV1(arcClient, d.Id(), nil, d.Get("tags"))
	if err != nil {
		return err
	}

	return resourceCCloudArcAgentV1Read(d, meta)
}

func resourceCCloudArcAgentV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	arcClient, err := config.arcV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack Arc client: %s", err)
	}

	agent, err := agents.Get(arcClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Unable to retrieve ccloud_arc_agent_v1")
	}

	arcCCloudArcAgentV1ReadAgent(d, arcClient, agent, GetRegion(d, config))

	return nil
}

func resourceCCloudArcAgentV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	arcClient, err := config.arcV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack Arc client: %s", err)
	}

	oldTags, newTags := d.GetChange("tags")
	err = updateArcAgentTagsV1(arcClient, d.Id(), oldTags, newTags)
	if err != nil {
		return err
	}

	return resourceCCloudArcAgentV1Read(d, meta)
}

func resourceCCloudArcAgentV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	arcClient, err := config.arcV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack Arc client: %s", err)
	}

	log.Printf("[DEBUG] Deleting ccloud_arc_agent_v1: %s", d.Id())
	err = agents.Delete(arcClient, d.Id()).ExtractErr()
	if err != nil {
		return CheckDeleted(d, err, "Error deleting ccloud_arc_agent_v1")
	}

	return nil
}
