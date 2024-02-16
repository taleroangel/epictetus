package storage

import "os"

// Where databases and images are stored
const StoragePath = "storage/"

// Create the storage path
func CreateStoragePath() error {
	// Create the directory
	err := os.Mkdir(StoragePath, os.ModePerm)
	// Check if it exists
	if os.IsExist(err) || err == nil {
		return nil
	}
	return err
}
