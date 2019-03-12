package ccloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/kayrus/gophercloud-arc/arc/v1/jobs"
)

type chefZeroPayload struct {
	OmnitruckUrl string                   `json:"omnitruck_url"`
	ChefVersion  string                   `json:"chef_version"`
	RunList      []string                 `json:"run_list"`
	RecipeURL    string                   `json:"recipe_url"`
	Attributes   map[string]interface{}   `json:"attributes"`
	Debug        bool                     `json:"debug"`
	Nodes        []map[string]interface{} `json:"nodes"`
	NodeName     string                   `json:"name"`
}

type tarballPayload struct {
	URL         string            `json:"url"`
	Path        string            `json:"path"`
	Arguments   []string          `json:"arguments"`
	Environment map[string]string `json:"environment"`
}

func arcCCloudArcJobV1BuildExecutePayload(v []interface{}) (action string, payload string) {
	for _, e := range v {
		execute := e.(map[string]interface{})

		if v, ok := execute["script"]; ok && len(v.(string)) > 0 {
			return "script", v.(string)
		}

		if v, ok := execute["tarball"]; ok {
			return "tarball", arcCCloudArcJobV1ParseTarball(v)
		}
	}

	return
}

func arcCCloudArcJobV1ParseTarball(v interface{}) string {
	var payload string

	for _, t := range v.([]interface{}) {
		var tarball tarballPayload
		tmp := t.(map[string]interface{})

		if val, ok := tmp["url"]; ok {
			tarball.URL = val.(string)
		}
		if val, ok := tmp["path"]; ok {
			tarball.Path = val.(string)
		}
		if val, ok := tmp["arguments"]; ok {
			tarball.Arguments = expandToStringSlice(val.([]interface{}))
		}
		if val, ok := tmp["environment"]; ok {
			tarball.Environment = expandToMapStringString(val.(map[string]interface{}))
		}

		bytes, _ := json.Marshal(tarball)
		payload = string(bytes[:])
	}

	return payload
}

func arcCCloudArcJobV1BuildChefPayload(v []interface{}) string {
	var payload string

	for _, c := range v {
		chef := c.(map[string]interface{})
		var chefZero chefZeroPayload

		if val, ok := chef["run_list"]; ok {
			chefZero.RunList = expandToStringSlice(val.([]interface{}))
		}
		if val, ok := chef["recipe_url"]; ok {
			chefZero.RecipeURL = val.(string)
		}
		if val, ok := chef["attributes"]; ok {
			json.Unmarshal([]byte(val.(string)), &chefZero.Attributes)
		}
		if val, ok := chef["debug"]; ok {
			chefZero.Debug = val.(bool)
		}
		if val, ok := chef["nodes"]; ok {
			json.Unmarshal([]byte(val.(string)), &chefZero.Nodes)
		}
		if val, ok := chef["node_name"]; ok {
			chefZero.NodeName = val.(string)
		}
		if val, ok := chef["omnitruck_url"]; ok {
			chefZero.OmnitruckUrl = val.(string)
		}
		if val, ok := chef["chef_version"]; ok {
			chefZero.ChefVersion = val.(string)
		}

		bytes, _ := json.Marshal(chefZero)
		payload = string(bytes[:])
	}

	return payload
}

func arcCCloudArcJobV1FlattenExecute(job *jobs.Job) ([]map[string]interface{}, error) {
	if !strSliceContains([]string{"tarball", "script"}, job.Action) {
		return []map[string]interface{}{}, nil
	}

	if job.Action == "script" {
		return []map[string]interface{}{{
			"script":  job.Payload,
			"tarball": []map[string]interface{}{},
		}}, nil
	}

	var tarball tarballPayload
	err := json.Unmarshal([]byte(job.Payload), &tarball)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tarball payload: %s", err)
	}

	return []map[string]interface{}{{
		"script": "",
		"tarball": []map[string]interface{}{{
			"url":         tarball.URL,
			"path":        tarball.Path,
			"arguments":   tarball.Arguments,
			"environment": tarball.Environment,
		}},
	}}, nil
}

func arcCCloudArcJobV1FlattenChef(job *jobs.Job) ([]map[string]interface{}, error) {
	if job.Action != "zero" {
		return []map[string]interface{}{}, nil
	}

	var chef chefZeroPayload

	err := json.Unmarshal([]byte(job.Payload), &chef)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal chef zero payload: %s", err)
	}

	attributes, _ := json.Marshal(chef.Attributes)
	nodes, _ := json.Marshal(chef.Nodes)

	return []map[string]interface{}{{
		"run_list":      chef.RunList,
		"recipe_url":    chef.RecipeURL,
		"attributes":    string(attributes[:]),
		"debug":         chef.Debug,
		"nodes":         string(nodes[:]),
		"node_name":     chef.NodeName,
		"omnitruck_url": chef.OmnitruckUrl,
		"chef_version":  chef.ChefVersion,
	}}, nil
}

func arcCCloudArcJobV1Filter(d *schema.ResourceData, arcClient *gophercloud.ServiceClient, resourceName string) ([]jobs.Job, error) {
	agentID := d.Get("agent_id").(string)
	timeout := d.Get("timeout").(int)
	agent := d.Get("agent").(string)
	action := d.Get("action").(string)
	status := d.Get("status").(string)

	listOpts := jobs.ListOpts{AgentID: agentID}

	log.Printf("[DEBUG] %s list options: %#v", resourceName, listOpts)

	allPages, err := jobs.List(arcClient, listOpts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("Unable to list %s: %s", resourceName, err)
	}

	allJobs, err := jobs.ExtractJobs(allPages)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve %s: %s", resourceName, err)
	}

	var jobs []jobs.Job
	for _, job := range allJobs {
		found := true
		if found && timeout > 0 && job.Timeout != timeout {
			found = false
		}
		if found && len(agent) > 0 && job.Agent != agent {
			found = false
		}
		if found && len(action) > 0 && job.Action != action {
			found = false
		}
		if found && len(status) > 0 && job.Status != status {
			found = false
		}

		if found {
			jobs = append(jobs, job)
		}
	}

	return jobs, nil
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
