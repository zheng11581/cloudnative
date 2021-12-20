package main

import (
	"fmt"
	"log"
)

type customer struct {
	name string
	sex  string
	age  int8
}

func (c *customer) hello() (*customer, error) {
	fmt.Printf("Hello, %s\n", c.name)
	err := fmt.Errorf("%s\n", "This is an error")
	return c, err
}

func (c *customer) bye() {
	fmt.Printf("Bye Bye, %s\n", c.name)
}

func main() {
	c := customer{
		name: "zhenghaicheng",
		sex:  "male",
		age:  33,
	}
	c1, err := c.hello()
	defer c1.bye()
	if err != nil {
		log.Fatal(err)
	}

	//res, err := http.Get("https://www.spreadsheetdb.io/")
	//defer res.Body.Close()
	//if err != nil {
	//	log.Fatal(err)
	//}
}
