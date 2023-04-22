package secrets

import "errors"

type MemoryVault struct {
	encodingKey string
	keyValues   map[string]string
}

func InMemory(encodingKey string) Vault {
	return Vault{
		encodingKey: encodingKey,
		keyValues:   make(map[string]string),
	}
}

func (v *MemoryVault) Get(key string) (string, error) {
	hex, ok := v.keyValues[key]
	if !ok {
		return "", errors.New("secret: no value for that key.")
	}
	ret, err := Decrypt(v.encodingKey, hex)
	if err != nil {
		return "", err
	}
	return ret, nil
}

func (v *MemoryVault) Set(key, value string) error {
	encryptedValue, err := Encrypt(v.encodingKey, value)
	if err != nil {
		return err
	}
	v.keyValues[key] = encryptedValue
	return nil
}
