package demoif

import (
	"fmt"
)

type IF interface {
	GetName() string
}

type Human struct {
	FirstName string
	LastName  string
}

type Plane struct {
	Vendor string
	Model  string
}

func (h Human) GetName() string {
	name := fmt.Sprintf("%s %s", h.FirstName, h.LastName)
	return name
}

func (p Plane) GetName() string {
	spec := fmt.Sprintf("%s %s", p.Vendor, p.Model)
	return spec

}

type Car struct {
	Factory string
	Model   string
}

func (c Car) GetName() string {
	name := fmt.Sprintf("%s %s", c.Factory, c.Model)
	return name
}
