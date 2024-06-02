provider "decrypt" {
  private_key = file("${path.module}/private.key")
}

resource "decrypt_file" "example" {
  file_path = "${path.module}/encrypted_file.yaml"
}

output "decrypted_content" {
  value = decrypt_file.example.decrypted_content
}
