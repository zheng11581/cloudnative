package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	content, err := os.ReadFile("C:/Users/ThinkPad/Desktop/phones1222.txt")
	if err != nil {
		panic(err)
	}
	phones := strings.Split(string(content), "\r")
	times := make(map[string]int)
	for _, phone := range phones {
		num, ok := times[phone]
		if !ok {
			times[phone] = 1
		}
		times[phone] = num + 1
	}
	for k, v := range times {
		fmt.Printf("phone: %s, times: %d\n", strings.Replace(k, "\n", "", -1), v)
	}

}
