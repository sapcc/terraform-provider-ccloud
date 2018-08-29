package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/sapcc/terraform-ccloud-provider/ccloud"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ccloud.Provider})
}
