package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"filippo.io/age"
	"gopkg.in/yaml.v3"
)

func resourceDecryptFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceDecryptFileCreate,
		Read:   resourceDecryptFileRead,
		Delete: resourceDecryptFileDelete,

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

func resourceDecryptFileCreate(d *schema.ResourceData, m interface{}) error {
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

func resourceDecryptFileRead(d *schema.ResourceData, m interface{}) error {
	// Read is not necessary for this example
	return nil
}

func resourceDecryptFileDelete(d *schema.ResourceData, m interface{}) error {
	// Delete is not necessary for this example
	return nil
}

func decryptFile(filePath, privateKey string) (string, error) {
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	identity, err := age.ParseX25519Identity(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	r := bytes.NewReader(fileData)
	dec, err := age.Decrypt(r, identity)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt file: %v", err)
	}

	decryptedData, err := ioutil.ReadAll(dec)
	if err != nil {
		return "", fmt.Errorf("failed to read decrypted data: %v", err)
	}

	// SOPS integration: Parse the decrypted content as a SOPS file.
	var sopsData map[string]interface{}
	err = yaml.Unmarshal(decryptedData, &sopsData)
	if err != nil {
		return "", fmt.Errorf("failed to parse decrypted YAML: %v", err)
	}

	// Check if the file is a SOPS file.
	if _, ok := sopsData["sops"]; ok {
		decryptedContent, err := parseSopsData(sopsData)
		if err != nil {
			return "", fmt.Errorf("failed to parse SOPS data: %v", err)
		}
		return decryptedContent, nil
	}

	// If not a SOPS file, return the raw decrypted content.
	return string(decryptedData), nil
}

func parseSopsData(sopsData map[string]interface{}) (string, error) {
	// Extract the encrypted values and decrypt them if necessary.
	// This example assumes the decrypted content is stored in a specific key.
	// Adjust the logic based on your SOPS file structure.
	decryptedContent, err := json.Marshal(sopsData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal decrypted data: %v", err)
	}

	return string(decryptedContent), nil
}
