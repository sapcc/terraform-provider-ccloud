---
layout: "ccloud"
page_title: "Converged Cloud: ccloud_automation_run_v1"
sidebar_current: "docs-ccloud-resource-automation-run-v1"
description: |-
  Run a Lyra Automation on the target.
---

# ccloud\_automation\_run\_v1

Use this resource to run a Lyra Automation on the target agent(s). The resource
will wait for the final Run state: `failed` or `completed`.

The `terraform destroy` command destroys the `ccloud_automation_run_v1` state,
but not the remote Lyra Automation Run object.

## Example Usage

```hcl
data "ccloud_automation_v1" "automation_1" {
  name = "chef-automation"
}

data "ccloud_arc_agent_v1" "agent_1" {
  filter  = "@metadata_name = 'hostname'"
  timeout = 600
}

resource "ccloud_automation_run_v1" "run_1" {
  automation_id = "${data.ccloud_automation_v1.automation_1.id}"
  selector      = "@identity = '${data.ccloud_arc_agent_v1.agent_1.id}'"
}

data "ccloud_arc_job_v1" "job" {
  job_id = "${ccloud_automation_run_v1.run_1.jobs[0]}"
}

output "run_log" {
  value = "${ccloud_automation_run_v1.run_1.log}"
}

output "job_log" {
  value = "${data.ccloud_arc_job_v1.job.log}"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the Automation client. If
  omitted, the `region` argument of the provider is used. Changing this forces
  a new resource to be created.

* `automation_id` - (Required) The ID of the Lyra automation to run. Changing
  this forces a new resource to be created.

* `selector` - (Required) The Arc Agent selector to run the Lyra automation.
  Changing this forces a new resource to be created.

## Attributes Reference

* `id` - The ID of the Lyra Automation Run.
* `region` - See Argument Reference above.
* `automation_id` - See Argument Reference above.
* `selector` - See Argument Reference above.
* `automation_name` - The name of the executed automation.
* `repository_revision` - The repository revision of the executed automation.
* `automation_attributes` - The attributes of the executed automation.
* `state` - The Automation Run state. Can either be `preparing`, `executing`,
  `failed` or `completed`.
* `log` - The Automation Run log.
* `jobs` - The list of the Arc Jobs ID, created by the Automation Run.
* `project_id` - The parent Openstack project ID.
* `created_at` - The date the Lyra automation was created.
* `updated_at` - The date the Lyra automation was last updated.
* `owner` - The user, who submitted the Automation Run. The structure is
  described below.

The `owner` attribute has fields below:

* `id` - The OpenStack user ID.

* `name` - The OpenStack user name.

* `domain_id` - The OpenStack domain ID.

* `domain_name` - The OpenStack domain name.
