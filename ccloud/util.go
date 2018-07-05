package ccloud

import "github.com/hashicorp/terraform/helper/schema"

func GetRegion(d *schema.ResourceData, config *Config) string {
	if v, ok := d.GetOk("region"); ok {
		return v.(string)
	}

	return config.Region
}
