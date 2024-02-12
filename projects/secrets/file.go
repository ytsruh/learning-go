package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

func File(encodingKey, filepath string) *FileVault {
	return &FileVault{
		encodingKey: encodingKey,
		filepath:    filepath,
	}
}

type FileVault struct {
	encodingKey string
	filepath    string
	mutex       sync.Mutex
	keyValues   map[string]string
}

func (v *FileVault) loadKeyValues() error {
	f, err := os.Open(v.filepath)
	if err != nil {
		v.keyValues = make(map[string]string)
		return nil
	}
	defer f.Close()
	var sb strings.Builder
	_, err = io.Copy(&sb, f)
	if err != nil {
		return err
	}
	decryptedJSON, err := Decrypt(v.encodingKey, sb.String())
	if err != nil {
		return err
	}
	r := strings.NewReader(decryptedJSON)
	dec := json.NewDecoder(r)
	err = dec.Decode(&v.keyValues)
	if err != nil {
		return err
	}
	return nil
}

func (v *FileVault) saveKeyValues() error {
	var sb strings.Builder // string builder has a io.writer
	enc := json.NewEncoder(&sb)
	err := enc.Encode(v.keyValues)
	if err != nil {
		return err
	}
	encryptedJSON, err := Encrypt(v.encodingKey, sb.String())
	if err != nil {
		return err
	}
	f, err := os.OpenFile(v.filepath, os.O_RDWR|os.O_CREATE, 0755) // taken from standard docs
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprint(f, encryptedJSON)
	if err != nil {
		return err
	}
	return nil
}

func (v *FileVault) Get(key string) (string, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	err := v.loadKeyValues()
	if err != nil {
		return "", err
	}
	value, ok := v.keyValues[key]
	if !ok {
		return "", errors.New("secret: no value for that key")
	}
	return value, nil
}

func (v *FileVault) Set(key, value string) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	err := v.loadKeyValues()
	if err != nil {
		return err
	}
	v.keyValues[key] = value
	err = v.saveKeyValues()
	return err
}
