package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"private_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Private key for decryption.",
				DefaultFunc: schema.EnvDefaultFunc("DECRYPT_PRIVATE_KEY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"decrypt_file": resourceDecryptFile(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"decrypt_file": dataSourceDecryptFile(),
		},
	}
}
