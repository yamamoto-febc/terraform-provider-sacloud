package sacloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_token": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"access_token_secret": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sacloud_server": resourceServer(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		AccessToken:       d.Get("access_token").(string),
		AccessTokenSecret: d.Get("access_token_secret").(string),
		Region:            d.Get("region").(string),
	}

	return config.NewClient(), nil
}
