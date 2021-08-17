---
layout: "ccloud"
page_title: "Converged Cloud: ccloud_domain_quota_v1"
sidebar_current: "docs-ccloud-datasource-domain-quota-v1"
description: |-
  Get information on the Domain Quota Resources
---

# ccloud\_domain\_quota\_v1

Use this data source to read the Limes (Quota) Domain resources.

## Example Usage

```hcl
data "openstack_identity_project_v3" "demo" {
  name = "demo"
}

data "ccloud_domain_quota_v1" "quota" {
  domain_id  = data.openstack_identity_project_v3.demo.domain_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the Limes client. If
  omitted, the `region` argument of the provider is used.

* `domain_id` – (Required) The ID of the domain to read the quota. Changing
  this forces a new resource to be created.

## Attributes Reference

In addition to arguments above, extra attributes are exported. Please refer
to the `ccloud_domain_quota_v1` resource arguments and attributes
[documentation](../resources/domain_quota_v1.html) for more information.
