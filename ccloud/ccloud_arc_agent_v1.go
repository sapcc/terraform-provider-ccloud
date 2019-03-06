package ccloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/kayrus/gophercloud-arc/arc/v1/agents"
	"github.com/kayrus/gophercloud-arc/arc/v1/jobs"
)

func arcCCloudArcAgentV1ReadAgent(d *schema.ResourceData, arcClient *gophercloud.ServiceClient, agent *agents.Agent, region string) {
	if len(agent.Facts) == 0 {
		facts, err := agents.GetFacts(arcClient, agent.AgentID).Extract()
		if err != nil {
			log.Printf("Unable to retrieve facts for ccloud_arc_agent_v1: %s", err)
		}
		agent.Facts = facts
		log.Printf("[DEBUG] Retrieved ccloud_arc_agent_v1 facts %s: %+v", agent.AgentID, agent.Facts)
	}

	d.Set("display_name", agent.DisplayName)
	d.Set("agent_id", agent.AgentID)
	d.Set("project", agent.Project)
	d.Set("organization", agent.Organization)
	d.Set("all_tags", agent.Tags)
	d.Set("created_at", agent.CreatedAt.Format(time.RFC3339))
	d.Set("updated_at", agent.UpdatedAt.Format(time.RFC3339))
	d.Set("updated_with", agent.UpdatedWith)
	d.Set("updated_by", agent.UpdatedBy)

	d.Set("facts", expandToMapStringString(agent.Facts))
	factsAgents, _ := agent.Facts["agents"]
	if v, ok := factsAgents.(map[string]interface{}); ok {
		d.Set("facts_agents", expandToMapStringString(v))
	} else {
		d.Set("facts_agents", map[string]string{})
	}

	d.Set("region", region)
}

func arcCCloudArcAgentV1GetAgent(arcClient *gophercloud.ServiceClient, agentID, filter string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var agent *agents.Agent
		var err error

		if len(agentID) == 0 && len(filter) == 0 {
			return nil, "", fmt.Errorf("At least one of agent_id or filter parameters is expected in ccloud_arc_agent_v1")
		}

		if len(agentID) > 0 {
			agent, err = agents.Get(arcClient, agentID).Extract()
			if err != nil {
				return nil, "", fmt.Errorf("Unable to retrieve %s ccloud_arc_agent_v1: %s", agentID, err)
			}
		} else {
			listOpts := agents.ListOpts{Filter: filter}

			log.Printf("[DEBUG] ccloud_arc_agent_v1 list options: %#v", listOpts)

			allPages, err := agents.List(arcClient, listOpts).AllPages()
			if err != nil {
				return nil, "", fmt.Errorf("Unable to list ccloud_arc_agent_v1: %s", err)
			}

			allAgents, err := agents.ExtractAgents(allPages)
			if err != nil {
				return nil, "", fmt.Errorf("Unable to retrieve ccloud_arc_agent_v1: %s", err)
			}

			if len(allAgents) == 0 {
				// Retryable case, when timeout is set
				return nil, "No ccloud_arc_agent_v1 found", nil
			}

			if len(allAgents) > 1 {
				return nil, "", fmt.Errorf("More than one ccloud_arc_agent_v1 found (%d)", len(allAgents))
			}

			agent = &allAgents[0]
		}

		log.Printf("[DEBUG] Retrieved ccloud_arc_agent_v1 %s: %+v", agent.AgentID, *agent)

		return agent, "active", nil
	}
}

func flattenArcJobUserV1(user jobs.User) []interface{} {
	return []interface{}{map[string]interface{}{
		"id":          user.ID,
		"name":        user.Name,
		"domain_id":   user.DomainID,
		"domain_name": user.DomainName,
		"roles":       user.Roles,
	}}
}

func updateArcAgentTagsV1(arcClient *gophercloud.ServiceClient, agentID string, oldTagsRaw, newTagsRaw interface{}) error {
	var tagsToDelete []string
	oldTags, _ := oldTagsRaw.(map[string]interface{})
	newTags, _ := newTagsRaw.(map[string]interface{})

	// Determine if any tag keys were removed from the configuration.
	// Then request those keys to be deleted.
	for oldKey := range oldTags {
		var found bool
		for newKey := range newTags {
			if oldKey == newKey {
				found = true
			}
		}

		if !found {
			tagsToDelete = append(tagsToDelete, oldKey)
		}
	}

	for _, key := range tagsToDelete {
		err := agents.DeleteTag(arcClient, agentID, key).ExtractErr()
		if err != nil {
			return fmt.Errorf("Error deleting %s tag from %s ccloud_arc_agent_v1: %s", key, agentID, err)
		}
	}

	// Update existing tags and add any new tags.
	tagsOpts := make(agents.Tags)
	for k, v := range newTags {
		tagsOpts[k] = v.(string)
	}

	err := agents.CreateTags(arcClient, agentID, tagsOpts).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error updating tags for %s ccloud_arc_agent_v1: %s", agentID, err)
	}

	return nil
}

func waitForArcJobV1(arcClient *gophercloud.ServiceClient, id string, target []string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for %s job to become %s.", id, target)

	stateConf := &resource.StateChangeConf{
		Target:     target,
		Pending:    pending,
		Refresh:    arcJobV1GetStatus(arcClient, id),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForState()

	return err
}

func arcJobV1GetStatus(arcClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		job, err := jobs.Get(arcClient, id).Extract()
		if err != nil {
			return nil, "", fmt.Errorf("Unable to retrieve %s ccloud_arc_job_v1: %s", id, err)
		}

		return job, job.Status, nil
	}
}
