package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/yamamoto-febc/terraform-provider-sacloud/sacloud"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: sacloud.Provider,
	})
}