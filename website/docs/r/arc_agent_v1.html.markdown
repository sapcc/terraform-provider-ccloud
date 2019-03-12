---
layout: "ccloud"
page_title: "Converged Cloud: ccloud_arc_agent_v1"
sidebar_current: "docs-ccloud-resource-arc-agent-v1"
description: |-
  Bring the Arc Agent resource under Terraform management
---

# ccloud\_arc\_agent\_v1

Use this resource to bring the existing Arc Agent resource under Terraform
management.

Unlike a regular Terraform resource, when it should be created by Terraform or
imported manually, the Arc Agent resource could not be created by Terraform
directly. The remote agent resource is being created by an Arc Agent, running
inside the OpenStack compute instance.

Nevertheless this Terraform resource allows to bring the Arc Agent under
Terraform management during the `create` stage, when it tries to find the
existing resource guided by the `filter` argument. It will wait for an Arc Agent
to be available (the Arc Agent bootstrap takes time due to compute instance boot
time and cloud-init execution delay) within the `timeouts`
[nested](/docs/configuration/resources.html#operation-timeouts) block argument.
The default create timeout is 10 minutes.

Once the Terraform manage the Arc Agent resource, Terraform can manage an agent
tags, and the `terraform destroy` command will destroy the remote Arc Agent
resource.

## Example Usage

### Manage an Arc Agent found with a filter

The example below will be completed once it finds the exact one agent
satisfying the specified filter.

```hcl
resource "ccloud_arc_agent_v1" "agent_1" {
  filter  = "@metadata_name = 'hostname'"

  timeouts = {
    create = "10m"
  }
}
```

### Manage an Arc Agent found with a filter and a compute instance dependency

The example below will be completed once it finds the exact one agent
satisfying the specified filter with an `openstack_compute_instance_v2` ID as a
dependency.

```hcl
resource "ccloud_arc_agent_bootstrap_v1" "agent_1" {}

resource "openstack_compute_instance_v2" "node" {
  name        = "linux-vm"
  image_name  = "ubuntu-16.04-amd64"
  flavor_name = "m1.small"
  user_data   = "${ccloud_arc_agent_bootstrap_v1.agent_1.user_data}"

  network = {
    name = "private_network"
  }
}

resource "ccloud_arc_agent_v1" "agent_1" {
  filter  = "@metadata_uuid = '${openstack_compute_instance_v2.node.id}'"

  timeouts = {
    create = "10m"
  }
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the Arc client. If
  omitted, the `region` argument of the provider is used. Changing this forces
  a new resource to be created.

* `filter` - (Required) The filter, used to filter the desired Arc Agent.
  Changing this forces a new resource to be created.

* `tags` - (Optional) The tags map to be appended to the Arc Agent. If an agent
  already has the tag key, specified as an argument, the key value will be
  overwritten to the value, defined in the resource.

## Attributes Reference

`id` is set to the ID of the found agent. In addition, the following attributes
are exported:

* `region` - See Argument Reference above.
* `filter` - See Argument Reference above.
* `tags` - See Argument Reference above.
* `display_name` - The Arc agent display name.
* `project` - The Arc agent parent OpenStack project ID.
* `organization` - The Arc agent parent OpenStack domain ID.
* `created_at` - The date the Arc agent was created.
* `updated_at` - The date the Arc agent was last updated.
* `updated_with` - The registration ID, used to submit the latest update.
* `updated_by` - The type of the application, submitted the latest update.
* `all_tags` - The map of tags, assigned on the Arc agent.
* `facts` - The map of facts, submitted by the Arc agent.
* `facts_agents` - The map of agent types enabled on the Arc agent.

## Import

An Arc Agent can be imported using the `id`, e.g.

```
$ terraform import ccloud_arc_agent_v1.agent_1 4107c3ea-0755-4a01-bfc4-cc4fe777ac98
```

The filter argument should be set to an empty value during the import.
