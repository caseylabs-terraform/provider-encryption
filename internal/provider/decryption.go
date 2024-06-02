package provider

import (
	"bytes"
	"errors"
	"io/ioutil"

	"filippo.io/age"
	"filippo.io/age/agessh"
	"gopkg.in/yaml.v2"
)

func decryptFile(privateKey, filePath string) (map[string]interface{}, error) {
	// Read the encrypted file
	encryptedData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Decrypt the file based on its format (sops or age)
	decryptedData, err := decryptWithAge(privateKey, encryptedData)
	if err != nil {
		return nil, err
	}

	// Convert the decrypted data to a map
	var decryptedValues map[string]interface{}
	err = yaml.Unmarshal(decryptedData, &decryptedValues)
	if err != nil {
		return nil, err
	}

	return decryptedValues, nil
}

func decryptWithAge(privateKey string, encryptedData []byte) ([]byte, error) {
	// Check if the private key is an SSH key
	identity, err := agessh.ParseIdentity([]byte(privateKey))
	if err != nil {
		return nil, errors.New("failed to parse private key")
	}

	// Decrypt using age
	r, err := age.Decrypt(bytes.NewReader(encryptedData), identity)
	if err != nil {
		return nil, err
	}

	decryptedData, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}
