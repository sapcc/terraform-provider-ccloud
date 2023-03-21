---
layout: "ccloud"
page_title: "Converged Cloud: ccloud_project_quota_v1"
sidebar_current: "docs-ccloud-datasource-project-quota-v1"
description: |-
  Get information on the Project Quota Resources
---

# ccloud\_project\_quota\_v1

Use this data source to read the Limes (Quota) Project resources.

## Example Usage

```hcl
data "openstack_identity_project_v3" "demo" {
  name = "demo"
}

data "ccloud_project_quota_v1" "quota" {
  domain_id  = data.openstack_identity_project_v3.demo.domain_id
  project_id = data.openstack_identity_project_v3.demo.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the Limes client. If
  omitted, the `region` argument of the provider is used.

* `domain_id` – (Required) The ID of the domain to read the quota. Changing
  this forces a new resource to be created.

* `project_id` - (Required) The ID of the project within the `domain_id` to
  read the quota. Changing this forces a new resource to be created.

## Attributes Reference

* `bursting` -  Contains information about the project bursting. The `bursting`
  object structure is documented below.

The `bursting` block supports:

* `enabled` - Indicates whether the quota bursting is enabled.

* `multiplier` - Indicates the quota bursting multiplier.

In addition to arguments above, extra attributes are exported. Please refer
to the `ccloud_project_quota_v1` resource arguments and attributes
[documentation](../resources/project_quota_v1.html) for more information.
