package main

import (
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"io"
	"log"
	"os"
)

type Product struct {
	Sku      int     `json:"sku" firestore:"sku"`
	Name     string  `json:"name" firestore:"name"`
	Type     string  `json:"type" firestore:"type"`
	Price    float64 `json:"price" firestore:"price"`
	Upc      string  `json:"upc" firestore:"upc"`
	Category []struct {
		ID   string `json:"id" firestore:"id"`
		Name string `json:"name" firestore:"name"`
	} `json:"category" firestore:"category"`
	Shipping     float64 `json:"shipping" firestore:"shipping"`
	Description  string  `json:"description" firestore:"description"`
	Manufacturer string  `json:"manufacturer" firestore:"manufacturer"`
	Model        string  `json:"model" firestore:"model"`
	URL          string  `json:"url" firestore:"url"`
	Image        string  `json:"image" firestore:"image"`
}

func main() {
	collection := os.Getenv("FIREBASE_COLLECTION")
	if collection == "" {
		log.Fatalln("No Firebase collection specified - do so with the environment variable FIREBASE_COLLECTION")
	}

	// Open products.json and serialise
	jsonFile, err := os.Open("products.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(jsonFile)
	byteValue, _ := io.ReadAll(jsonFile)
	var products []Product
	err = json.Unmarshal(byteValue, &products)
	if err != nil {
		log.Fatalln(err)
	}
	// Connect to Firestore
	ctx := context.Background()
	fb := &firebase.Config{}
	app, err := firebase.NewApp(ctx, fb)
	if err != nil {
		log.Fatalln(err)
	}
	fs, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Inserting %v products to Firestore\n", len(products))

	bw := fs.BulkWriter(ctx)

	for _, product := range products {
		_, err := bw.Set(fs.Collection(collection).Doc(fmt.Sprintf("%v", product.Sku)), product)
		if err != nil {
			log.Printf("Error inserting %v: %v\n", product.Sku, err)
		}
	}
	bw.End()
	bw.Flush()
	log.Println("Done")
}
