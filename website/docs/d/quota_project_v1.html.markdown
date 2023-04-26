---
layout: "ccloud"
page_title: "Converged Cloud: ccloud_quota_project_v1"
sidebar_current: "docs-ccloud-datasource-quota-project-v1"
description: |-
  Get information on the Project Quota Resources
---

# ccloud\_quota\_project\_v1

Use this data source to read the Limes (Quota) Project resources.

## Example Usage

```hcl
data "ccloud_quota_project_v1" "quota" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the Limes client. If
  omitted, the `region` argument of the provider is used.

* `domain_id` – (Optional) The ID of the domain to read the quota. Defaults to
  the current domain scope. Changing this forces a new resource to be created.

* `project_id` - (Optional) The ID of the project within the `domain_id` to
  read the quota. Defaults to the current project scope. Changing this forces a
  new resource to be created.

## Attributes Reference

* `bursting` -  Contains information about the project bursting. The `bursting`
  object structure is documented below.

The `bursting` block supports:

* `enabled` - Indicates whether the quota bursting is enabled.

* `multiplier` - Indicates the quota bursting multiplier.

In addition to arguments above, extra attributes are exported. Please refer
to the `ccloud_quota_project_v1` resource arguments and attributes
[documentation](../resources/quota_project_v1.html) for more information.
