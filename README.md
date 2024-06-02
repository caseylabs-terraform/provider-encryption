# provider-encryption

Terraform Encryption Provider (usings `sops` and `age`).

- **Note:** this provider is used for internal projects and is not publicly supported.

Status: `in development`

A Terraform provider that will:

- Decrypt an input file using `sops` and `age`.

- The file location can be local, or in a remote git repo.

- The provider will use a private key retrieved from either a local file, or from the value of an environment variable.

- The provider will decrypt the file using either an `age` private key, or an SSH private key, or an AWS KMS key.

- The provider will determine if an `age` or SSH private key has been provided. If it is an SSH key, the provider will temporarily convert it to an `age` to perform decryption.

- The provider will then be able to retrieve values from the encrypted file and expose them as Terraform values that can be used.