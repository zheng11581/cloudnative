package main

import (
	"cncamp/ch1/onclass/interface"
	"encoding/json"
	"fmt"
)

func Unmarshal2SHuman(objStr string) demoif.IF {

	h := demoif.Human{}
	err := json.Unmarshal([]byte(objStr), &h)
	if err != nil {
		fmt.Println(err)
	}
	return h
}

func Marshal2String(obj *demoif.IF) string {
	objBytes, err := json.Marshal(&obj)
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
		objString := Marshal2String(&f)
		obj := Unmarshal2SHuman(objString)
		fmt.Println(obj) // 只有Human被Unmarshal
	}

}
