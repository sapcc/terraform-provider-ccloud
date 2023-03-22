---
subcategory: "Guides"
layout: "ccloud"
page_title: "Migrating to a renamed resource"
sidebar_current: "docs-ccloud-migrating-between-renamed-resources"
description: |-
  This page documents how to migrate between two resources in the Converged Cloud Provider which have been renamed.
---

# Migrating to a renamed resource

This guide shows how to migrate from a resource which have been deprecated to
its replacement. The complete list of resources which have been deprecated in
favour of others can be found below.

!> **Note:** The following resources have been Deprecated and will be removed
in version 1.6 of the Converged Cloud Provider

| Old Name                              | New Name                |
| --------------------------------------| ------------------------|
| ccloud_quota                          | ccloud_quota_project_v1 |
| ccloud_quota_v1                       | ccloud_quota_project_v1 |
| ccloud_project_quota_v1               | ccloud_quota_project_v1 |
| ccloud_domain_quota_v1                | ccloud_quota_domain_v1  |
| ccloud_kubernetes                     | ccloud_kubernetes_v1    |
| ccloud_domain_quota_v1 (Data Source)  | ccloud_quota_domain_v1  |
| ccloud_project_quota_v1 (Data Source) | ccloud_quota_project_v1 |

## Migrating to a renamed resource

Rename the resource definition in your Terraform Configuration, e.g.

```diff
-resource "ccloud_project_quota_v1" "quota" {
+resource "ccloud_quota_project_v1" "quota" {
```

Then import the new resource by running:

```sh
$ terraform import ccloud_quota_project_v1.quota DOMAINID/PROJECTID
```

Domain and Project IDs can be retrived analyzying an old resource state:

```sh
$ terraform state show ccloud_project_quota_v1.quota | egrep 'project_id|domain_id'
```

Now we can remove an old resource state using `terraform state rm`, e.g.

```sh
$ terraform state rm ccloud_project_quota_v1.quota
```
