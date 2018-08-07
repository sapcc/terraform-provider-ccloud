Converged Cloud Enterprise Edition - Terraform Provider
=======================================================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Maintainers
-----------

This provider plugin is maintained by:

  * Michael Schmidt (@bugroger)

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)

Usage
---------------------


The CCloudEE provider is an extension to the [OpenStack Terraform
Provider](https://github.com/terraform-providers/terraform-provider-openstack).
It provides resources that allow to use Terraform for Converged Cloud's
additional services:

  * Limes for Quota Management
  * Kubernikus (Kubernetes as a Service)

The provider needs to be configured with the proper OpenStack credentials
before it can be used. For details see the OpenStack provider.


```
provider "ccloud" {
  auth_url         = "${var.auth_url}"
  region           = "${var.region}"
  user_name        = "${var.user_name}"
  user_domain_name = "${var.user_domain_name}"
  password         = "${var.password}"
  tenant_name      = "${var.tenant_name}"
  domain_name      = "${var.domain_name}"
}

data "openstack_identity_project_v3" "demo" {
  name = "${var.tenant_name}"
}

resource "ccloud_quota" "quota" {
  domain_id  = "${openstack_identity_project_v3.demo.domain_id}"
  project_id = "${openstack_identity_project_v3.demo.id}"

  compute {
    instances = 8
    cores     = 32
    ram       = 81920 
  }

  volumev2 {
    capacity  = 1024
    snapshots = 0
    volumes   = 128
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

  object-store {
    capacity = 1073741824
  }
}

resource "ccloud_kubernetes" "demo" {
  name           = "demo"
  ssh_public_key = "ssh-rsa AAAABHTmDMP6w=="

  node_pools = [
    { name = "payload0", flavor = "m1.xlarge_cpu", size = 2 },
    { name = "payload1", flavor = "m1.xlarge_cpu", size = 1 }
  ]
}

```

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/sapcc/terraform-ccloud-provider`

```sh
$ mkdir -p $GOPATH/src/github.com/sapcc; cd $GOPATH/src/github.com/sapcc
$ git clone git@github.com:sapcc/terraform-ccloud-provider
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/sapcc/terraform-ccloud-provider
$ make build
```


Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.10+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-ccloud
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```
