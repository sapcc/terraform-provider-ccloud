package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/sapcc/terraform-provider-ccloud/ccloud"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ccloud.Provider})
}
