package json

import (
	"encoding/json"
	"fmt"
	"os"
)

func Read[T any](path string) (data []T, err error) {

	dir := fmt.Sprintf("./../../data/%s", path)
	rawData, err := os.ReadFile(dir)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	err = json.Unmarshal(rawData, &data)
	if err != nil {
		fmt.Println("Error Unmarshal file:", err)
		return
	}

	return
}
