package provider

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestProvider(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: map[string]*terraform.ResourceProvider{
			"encryption": Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCheckEncryptionDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDecryptedValues("data.encryption_decrypt.test"),
				),
			},
		},
	})
}

const testAccCheckEncryptionDataSourceConfig = `
provider "encryption" {
  private_key = "example_private_key"
}

data "encryption_decrypt" "test" {
  file_path = "example_encrypted_file_path"
}
`

func testAccCheckDecryptedValues(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return errors.New("not found: " + n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		// Check the decrypted values here
		return nil
	}
}
