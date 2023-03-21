---
layout: "ccloud"
page_title: "Converged Cloud: ccloud_project_quota_v1"
sidebar_current: "docs-ccloud-resource-project-quota-v1"
description: |-
  Manages Project Quota Resources
---

# ccloud\_project\_quota\_v1

Manages Limes (Quota) Project resources.

~> **Note:** The `terraform destroy` command destroys the
`ccloud_project_quota_v1` state, but not the actual Limes project quota.

~> **Note:** You _must_ have admin privileges in your OpenStack cloud to use
this resource.

## Example Usage

```hcl
data "openstack_identity_project_v3" "demo" {
  name = "demo"
}

resource "ccloud_project_quota_v1" "quota" {
  domain_id  = data.openstack_identity_project_v3.demo.domain_id
  project_id = data.openstack_identity_project_v3.demo.id

  bursting {
    enabled = true
  }

  compute {
    instances = 8
    cores     = 32
    ram       = 81920
  }

  volumev2 {
    capacity               = 1024
    snapshots              = 0
    volumes                = 128
    capacity_standard_hdd  = 4096
    snapshots_standard_hdd = 0
    volumes_standard_hdd   = 1000
  }

  network {
    floating_ips         = 4
    networks             = 1
    ports                = 512
    routers              = 2
    security_group_rules = 64
    security_groups      = 4
    subnets              = 1
    healthmonitors       = 10
    l7policies           = 8
    listeners            = 16
    loadbalancers        = 8
    pools                = 8
    pool_members         = 10
    bgpvpns              = 0
    trunks               = 0
  }

  dns {
    zones      = 1
    recordsets = 16
  }

  sharev2 {
    share_networks    = 1
    share_capacity    = 1024
    shares            = 16
    snapshot_capacity = 512
    share_snapshots   = 8
  }

  objectstore {
    capacity = 1073741824
  }

  keppel {
    images = 100
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the Limes client. If
  omitted, the `region` argument of the provider is used. Changing this forces
  a new resource to be created.

* `domain_id` – (Required) The ID of the domain to manage the quota. Changing
  this forces a new resource to be created.

* `project_id` - (Required) The ID of the project within the `domain_id` to
  manage the quota. Changing this forces a new resource to be created.

* `bursting` - (Optional) If bursting is enabled, the project may exceed its
  granted quota by a certain multiplier. In case the higher usage level becomes
  permanent, users should request a quota extension from their domain admin or
  cloud admin. This is because burst usage is typically billed at a higher price
  than regular usage. The `bursting` object structure is documented below.

* `compute` - (Optional) The list of compute resources quota. Consists of
  `cores`, `instances`, `ram` (Mebibytes), `server_groups` and
  `server_group_members`.

* `volumev2` - (Optional) The list of block storage resources quota. Consists of
  `capacity` (Gibibytes), `snapshots` and `volumes`.

* `network` - (Optional) The list of network resources quota. Consists of
  `floating_ips`, `networks`, `ports`, `rbac_policies`, `routers`,
  `security_group_rules`, `security_groups`, `subnet_pools`, `subnets`,
  `healthmonitors`, `l7policies`, `listeners`, `loadbalancers`, `pools`,
  `pool_members`, `bgpvpns` and `trunks`.

* `dns` - (Optional) The list of DNS resources quota. Consists of `zones` and
  `recordsets`.

* `sharev2` - (Optional) The list of Shared File Systems resources quota.
  Consists of `share_networks`, `share_capacity` (Gibibytes), `shares`,
  `snapshot_capacity` (Gibibytes) and `share_snapshots`.

* `objectstore` - (Optional) The list of Object Storage resources quota.
  Consists of `capacity` (Bytes).

* `keppel` - (Optional) The list of Image Registry resources quota. Consists of
  `images`.

The `bursting` block supports:

* `enabled` - (Optional) Enables or disables the quota bursting.

* `multiplier` - (Computed) Indicates the quota bursting multiplier.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The project ID.

## Import

Limes Project Quota can be imported using the `domain_id` and `project_id`
arguments, e.g.

```
$ terraform import ccloud_project_quota_v1.demo bf2273b5-2926-4495-9fb7-f28c3abed5f6/ec407270-0249-4a82-a331-90ede2e78d9c
```
