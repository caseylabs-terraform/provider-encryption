package provider

import (
	"context"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"private_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The private key used to decrypt the file.",
				Sensitive:   true,
			},
			"private_key_env": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The environment variable containing the private key or the path to the private key file.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"encryption_decrypt": dataSourceDecrypt(),
		},
	}
}

func dataSourceDecrypt() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDecryptRead,
		Schema: map[string]*schema.Schema{
			"file_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The path to the encrypted file.",
			},
			"decrypted_values": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The decrypted values from the file.",
			},
		},
	}
}

func dataSourceDecryptRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var privateKey string
	if v, ok := d.GetOk("private_key"); ok {
		privateKey = v.(string)
	}
	privateKeyEnv := d.Get("private_key_env").(string)

	if privateKey == "" && privateKeyEnv != "" {
		privateKeyValue := os.Getenv(privateKeyEnv)
		if privateKeyValue == "" {
			return diag.Errorf("Environment variable %s is not set", privateKeyEnv)
		}

		// Check if the privateKeyValue is a path to a file
		if strings.HasPrefix(privateKeyValue, "/") || strings.HasPrefix(privateKeyValue, "./") {
			keyData, err := os.ReadFile(privateKeyValue)
			if err != nil {
				return diag.FromErr(err)
			}
			privateKey = string(keyData)
		} else {
			privateKey = privateKeyValue
		}
	}

	if privateKey == "" {
		return diag.Errorf("No private key provided")
	}

	filePath, ok := d.GetOk("file_path")
	if !ok {
		return diag.Errorf("file_path is required")
	}

	decryptedValues, err := decryptFile(privateKey, filePath.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("decrypted_values", decryptedValues); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(filePath.(string))

	return nil
}

// Mocked decryptFile function for demonstration purposes
func decryptFile(privateKey, filePath string) (map[string]interface{}, error) {
	// Implement your decryption logic here
	return map[string]interface{}{
		"example_key": "example_value",
	}, nil
}
