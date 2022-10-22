package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Product struct {
	Sku      int     `json:"sku"`
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Price    float64 `json:"price"`
	Upc      string  `json:"upc"`
	Category []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
	Shipping     float64 `json:"shipping"`
	Description  string  `json:"description"`
	Manufacturer string  `json:"manufacturer"`
	Model        string  `json:"model"`
	URL          string  `json:"url"`
	Image        string  `json:"image"`
}

func main() {
	// Open products.json and serialise
	jsonFile, err := os.Open("products.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(jsonFile)
	byteValue, _ := io.ReadAll(jsonFile)
	var products []Product
	err = json.Unmarshal(byteValue, &products)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < 10; i++ {
		fmt.Printf("%v = %v\n", products[i].Sku, products[i].Name)
	}
}
