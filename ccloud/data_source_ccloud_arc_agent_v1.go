package ccloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/kayrus/gophercloud-arc/arc/v1/agents"
)

func dataSourceCCloudArcAgentV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCCloudArcAgentV1Read,

		// Terraform timeouts don't work in data sources.
		// However "Timeouts" has to be specified, otherwise "timeouts" argument below won't work.
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(0),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"agent_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"filter"},
				ValidateFunc:  validation.NoZeroValues,
			},

			"filter": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"agent_id"},
				ValidateFunc:  validation.NoZeroValues,
			},

			// Terraform timeouts don't work in data sources.
			// This is a workaround.
			"timeouts": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"read": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateTimeout,
						},
					},
				},
			},

			// computed attributes
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

func dataSourceCCloudArcAgentV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	arcClient, err := config.arcV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack Arc client: %s", err)
	}

	var tmp interface{}
	var msg string
	agentID := d.Get("agent_id").(string)
	filter := d.Get("filter").(string)

	timeout, err := arcAgentV1ParseTimeout(d.Get("timeouts"))
	if err != nil {
		return fmt.Errorf("Error parsing the read timeout for ccloud_arc_job_v1: %s", err)
	}

	if timeout > 0 {
		// Retryable case, when timeout is set
		waitForAgent := &resource.StateChangeConf{
			Target:     []string{"active"},
			Refresh:    arcCCloudArcAgentV1GetAgent(arcClient, agentID, filter),
			Timeout:    timeout,
			Delay:      1 * time.Second,
			MinTimeout: 1 * time.Second,
		}
		tmp, err = waitForAgent.WaitForState()
	} else {
		// When timeout is not set, just get the agent
		tmp, msg, err = arcCCloudArcAgentV1GetAgent(arcClient, agentID, filter)()
	}

	if len(msg) > 0 && msg != "active" {
		return fmt.Errorf(msg)
	}

	if err != nil {
		return err
	}

	agent := tmp.(*agents.Agent)

	d.SetId(agent.AgentID)

	arcCCloudArcAgentV1ReadAgent(d, arcClient, agent, GetRegion(d, config))

	return nil
}
