package jsonUtils

import (
	"fmt"
	"os"
)

func getBaseDir() (dir string, err error) {

	// Get current working directory to build absolute path
	dir, err = os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	return
}
