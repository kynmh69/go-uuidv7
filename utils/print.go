package utils

import "fmt"

func PrintUUIDs(uuids []string) {
	for _, uuid := range uuids {
		fmt.Println(uuid)
	}
}