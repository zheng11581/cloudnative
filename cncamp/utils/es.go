package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

func DeleteIndex(index map[string]string) {
	req, err := http.NewRequest(http.MethodDelete, "http://es.example.com:31114/"+index["index"], nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deleting index %s, result is %s\n", index["index"], string(result))
}

func GetIndex() []map[string]string {
	req, err := http.NewRequest(http.MethodGet, "http://es.example.com:31114/_cat/indices?format=json&pretty", nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%v\n", string(result))

	var indexes []map[string]string
	err = json.Unmarshal(result, &indexes)
	if err != nil {
		panic(err)
	}

	reg1, err := regexp.Compile(".*kibana_1.*")
	if err != nil {
		panic(err)
	}

	for i, index := range indexes {
		value, ok := index["index"]
		if !ok {
			panic("No index")
		}
		if reg1.Match([]byte(value)) {
			indexes = append(indexes[:i], indexes[i+1:]...)
		}
		//if value == ".kibana_task_manager" || value == ".kibana_1\n" {
		//	indexes = append(indexes[:i], indexes[i+1:]...)
		//}
		//fmt.Printf("%v\n", index["index"])
	}

	reg2, err := regexp.Compile(".*kibana_task_manager.*")
	if err != nil {
		panic(err)
	}
	for i, index := range indexes {
		value, ok := index["index"]
		if !ok {
			panic("No index")
		}
		if reg2.Match([]byte(value)) {
			indexes = append(indexes[:i], indexes[i+1:]...)
		}
	}

	return indexes
}
