package usecase

import (
	"fmt"
	"reflect"
)

var count = 1

func GetTrackInfoBatch[I any, O any](url string, identifies []I) (data []O, err error) {

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

func GetTrackInfoBatchGetID[I any, O any](url string, identifies []I, getID func(I) int) (data []O, err error) {

	for _, ident := range identifies {

		id := getID(ident)

		oneData, err := GetTrackInfo[O](url, id)
		if err != nil {
			fmt.Printf("Error Getting Track info: %v", err)
			return data, err
		}

		data = append(data, oneData)
		fmt.Printf("%d - %d :\n", count, id)
		count++
	}

	return
}
