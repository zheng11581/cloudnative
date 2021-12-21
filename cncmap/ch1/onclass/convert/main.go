package main

import (
	demoif "cncamp/ch1/onclass/interface"
	"encoding/json"
	"fmt"
	"strconv"
)

func main() {
	// 类似类型间转换
	var x int8 = 10
	var y int32
	var z int64
	v := int32(x) + y
	w := int64(x) + z
	fmt.Printf("w=%d, v=%d\n", w, v)

	bs := []byte("hello world")
	str := string(bs)
	fmt.Printf("bs=%s, str=%s\n", bs, str)

	// 数字--->字符串转换
	strconv.FormatInt(int64(v), 10)
	strconv.FormatInt(w, 10)

	abc, _ := strconv.ParseInt("123", 10, 32)
	fmt.Printf("abc=%v\n", abc)

	hin := demoif.Human{
		FirstName: "王",
		LastName:  "力宏",
	}
	hhBytes, _ := json.Marshal(hin)
	fmt.Printf("json: %v\n", string(hhBytes))
	//hout := demoif.Human{} // 仅能转换为Human
	var hout interface{} // 可以转换为任何类型
	hout = demoif.Human{}
	_ = json.Unmarshal(hhBytes, &hout)
	houtMap, _ := hout.(map[string]interface{})
	for k, v := range houtMap {
		switch value := v.(type) {
		case interface{}:
			fmt.Printf("%s: %s\n", k, value)
		}
	}
}
