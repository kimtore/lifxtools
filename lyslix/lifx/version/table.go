package version

import (
	"encoding/json"
	"fmt"
)

type Vendor struct {
	ID       int `json:"vid"`
	Products []Product
}

type Product struct {
	ID       int `json:"pid"`
	Name     string
	Features Features
}

type Features struct {
	HEV               bool `json:"hev"`
	Color             bool
	Chain             bool
	Matrix            bool
	Relays            bool
	Buttons           bool
	Infrared          bool
	Multizone         bool
	TemperatureRange  []int `json:"temperature_range"`
	ExtendedMultizone bool  `json:"extended_multizone"`
}

type key struct {
	Vendor  int
	Product int
}

var products = make(map[key]*Product)

func init() {
	table := make([]Vendor, 0)
	err := json.Unmarshal([]byte(dataTable), &table)
	if err != nil {
		panic(err)
	}
	for _, vendor := range table {
		for i, product := range vendor.Products {
			k := key{vendor.ID, product.ID}
			products[k] = &vendor.Products[i]
		}
	}
}

func Lookup(vendor int, product int) (*Product, error) {
	p := products[key{vendor, product}]
	if p == nil {
		return nil, fmt.Errorf("product not found in database")
	}
	return p, nil
}
