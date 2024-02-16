package storage

import "os"

// Where databases and images are stored
const StoragePath = "storage/"

// Create the storage path
func CreateStoragePath() error {
	err := os.Mkdir(StoragePath, 0777)
	if err != nil && os.IsExist(err) {
		return nil
	}

	return err
}
