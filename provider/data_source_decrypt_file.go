package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDecryptFile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDecryptFileRead,

		Schema: map[string]*schema.Schema{
			"file_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Path to the encrypted file.",
			},
			"decrypted_content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Decrypted content of the file.",
			},
			"private_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Private key for decryption.",
				DefaultFunc: schema.EnvDefaultFunc("DECRYPT_PRIVATE_KEY", nil),
			},
		},
	}
}

func dataSourceDecryptFileRead(d *schema.ResourceData, m interface{}) error {
	privateKey := d.Get("private_key").(string)
	filePath := d.Get("file_path").(string)

	decryptedContent, err := decryptFile(filePath, privateKey)
	if err != nil {
		return err
	}

	d.Set("decrypted_content", decryptedContent)
	d.SetId(filePath) // Using file_path as ID
	return nil
}
