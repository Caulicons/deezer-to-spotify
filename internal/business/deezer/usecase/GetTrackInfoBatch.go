package usecase

import (
	"fmt"
	"reflect"
)

// This is a options one  to get the ID
func GetTrackInfoBatch[I any, O any](url string, identifies []I) (data []O, err error) {
	var count = 1

	for _, ident := range identifies {
		if count == 3 {
			return
		}

		v := reflect.ValueOf(ident)
		id := v.FieldByName("ID").Interface()

		fmt.Println("id :", id)
		count++
	}

	return
}

// This is other way to get the ID, i like this
func GetTrackInfoBatchGetID[I any, O any](url string, identifies []I, getID func(I) int, getTitle func(I) string) (data []O, err error) {
	var count = 1

	fmt.Println("🚩 Start Getting Deezer Track Info's: ")
	for _, ident := range identifies {

		id := getID(ident)

		oneData, err := GetTrackInfo[O](url, id)
		if err != nil {
			fmt.Printf("Error Getting Track info: %v", err)
			return data, err
		}

		data = append(data, oneData)
		fmt.Printf("%d - %s \n", count, getTitle(ident))
		count++
	}

	return
}
