package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	status := [...]string{"I", "am", "stupid", "and", "weak"}
	fmt.Printf("%v\n", status)
	for i := 0; i < len(status); i++ {
		switch i {
		case 2:
			status[i] = "smart"
		case 4:
			status[i] = "strong"
		default:
			continue
		}
	}
	fmt.Printf("%v\n", status)
	statusStr := Marshal2String(status)
	fmt.Printf("%v\n", statusStr)

}

func Marshal2String(obj interface{}) string {
	objBytes, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return string(objBytes)
}
