---
layout: "ccloud"
page_title: "Converged Cloud: ccloud_quota_domain_v1"
sidebar_current: "docs-ccloud-datasource-quota-domain-v1"
description: |-
  Get information on the Domain Quota Resources
---

# ccloud\_quota\_domain\_v1

Use this data source to read the Limes (Quota) Domain resources.

## Example Usage

```hcl
data "ccloud_quota_domain_v1" "quota" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the Limes client. If
  omitted, the `region` argument of the provider is used.

* `domain_id` – (Optional) The ID of the domain to read the quota. Defaults to
  the current domain scope. Changing this forces a new resource to be created.

## Attributes Reference

In addition to arguments above, extra attributes are exported. Please refer
to the `ccloud_quota_domain_v1` resource arguments and attributes
[documentation](../resources/quota_domain_v1.html) for more information.
