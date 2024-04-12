package storage

import (
	"encoding/json"
	"errors"
	"log"
	"os"
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

// Delete deletes the value for a given key.
func (kv *KeyValueStorage) Delete(key string) error {
	if _, ok := kv.data[key]; !ok {
		return errors.New("key does not exist")
	}
	delete(kv.data, key)
	kv.SaveToDisk()
	return nil
}
