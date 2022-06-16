package main

import (
	"encoding/json"
	"fmt"
	"log"

	demoif "github.com/zheng11581/cloudnative/cncamp/ch1/onclass/interface"
)

func Unmarshal2Struct(objStr string) interface{} {
	var obj interface{}
	err := json.Unmarshal([]byte(objStr), &obj)
	if err != nil {
		log.Fatal(err)
	}
	objMap, ok := obj.(map[string]interface{}) // map[string]interface{}和[]interface{}保存任意对象
	if !ok {
		log.Fatal("panic")
	}
	for k, v := range objMap {
		switch value := v.(type) {
		case interface{}:
			return fmt.Sprintf("%s %s", k, value) // 根据value类型判断
			//case string:
			//	return fmt.Sprintf("%s %s", k, value)
		}
	}
	return obj
}

func Marshal2String(obj interface{}) string {
	objBytes, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
	}
	return string(objBytes)
}

func main() {
	var interfaces []demoif.IF
	h := demoif.Human{
		FirstName: "zheng",
		LastName:  "haicheng",
	}
	interfaces = append(interfaces, h)

	p := demoif.Plane{
		Vendor: "XSpace",
		Model:  "new",
	}
	interfaces = append(interfaces, p)

	c := new(demoif.Car)
	c.Factory = "Tesla"
	c.Model = "Y"
	interfaces = append(interfaces, c)

	for _, f := range interfaces {
		//fmt.Println(f.GetName())
		objString := Marshal2String(f)
		obj := Unmarshal2Struct(objString)
		fmt.Println(obj)
	}
}
