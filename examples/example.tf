terraform {
  required_providers {
    encryption = {
      source  = "github.com/caseylabs-terraform/encryption"
    }
  }
}

provider "encryption" {
  private_key = "example_private_key"
}

data "encryption_decrypt" "test" {
  file_path = "example_encrypted_file_path"
}

output "decrypted_values" {
  value = data.encryption_decrypt.test.decrypted_values
}
