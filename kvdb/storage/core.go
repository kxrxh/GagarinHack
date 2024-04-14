package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type KeyValueStorage struct {
	data     map[string]string
	filePath string
}

var KeyValueStore *KeyValueStorage

// NewKeyValueStore creates a new instance of KeyValueStore.
func NewKeyValueStore(filePath string) error {
	KeyValueStore = &KeyValueStorage{
		data:     make(map[string]string),
		filePath: filePath,
	}
	err := KeyValueStore.loadFromDisk()
	return err
}

// loadFromDisk loads the data from disk on startup.
func (kv *KeyValueStorage) loadFromDisk() error {
	// Check if the file exists
	if _, err := os.Stat(kv.filePath); os.IsNotExist(err) {
		log.Println("Creating new empty storage...")
		return nil
	}

	// Read the file
	file, err := os.ReadFile(kv.filePath)
	if err != nil {
		return errors.New("unable to read file")
	}

	// Unmarshal the data into the map
	if err := json.Unmarshal(file, &kv.data); err != nil {
		log.Println(kv.filePath, file)
		return errors.New("unable to unmarshal file")
	}

	return nil
}

// SaveToDisk saves the data to disk after each transaction.
func (kv *KeyValueStorage) SaveToDisk() error {
	// Marshal the data
	data, err := json.Marshal(kv.data)
	if err != nil {
		return err
	}

	// Write the data to the file
	if err := os.WriteFile(kv.filePath, data, 0644); err != nil {
		return err
	}

	return nil
}

// Get retrieves the value for a given key.
func (kv *KeyValueStorage) Get(key string) (string, bool) {

	value, ok := kv.data[key]
	return value, ok
}

// Set sets the value for a given key.
func (kv *KeyValueStorage) Set(key, value string) {
	kv.data[key] = value
	kv.SaveToDisk()
}

// Delete deletes the value for a given key and all its variants.
func (kv *KeyValueStorage) Delete(key string) error {
	deleted := kv.DeleteSubtree(key)
	
	if _, ok := kv.data[key]; !ok {
		if deleted == nil {
			return nil
		}
		return errors.New("key does not exist")
	}
	delete(kv.data, key)

	kv.SaveToDisk()

	return nil
}

func (kv *KeyValueStorage) DeleteSubtree(key string) error {
	deleted := false

	// Iterate over the keys in kv.data
	for k := range kv.data {
		// Check if the key starts with the specified prefix
		if strings.HasPrefix(k, fmt.Sprintf("%s.", key)) {
			// Delete the key
			delete(kv.data, k)
			deleted = true
		}
	}
	// If no keys were deleted, return an error
	if !deleted {
		return errors.New("key or its variants do not exist")
	}

	// Save the updated data to disk
	kv.SaveToDisk()

	return nil
}
