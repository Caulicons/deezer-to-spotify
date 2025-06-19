package json

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// For agreement (by my self only) all the json file are storage in the folder
func Write[T any](data T, path string) error {

	// Create data directory if it doesn't exist
	dataDir := fmt.Sprintf("./../../data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Error creating data directory: %v", err)
		return err
	}

	destPath := fmt.Sprintf("%s/%s", dataDir, path)
	file, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Error opening the file %s: %v", destPath, err)
		return err
	}
	defer file.Close()

	jsonData, nil := json.Marshal(data)

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("Error writing file %s: %v", destPath, err)
		return err
	}

	return nil
}
