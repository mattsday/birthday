package main

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"google.golang.org/api/iterator"
	"log"
	"math/rand"
	"os"
)

var fs *firestore.Client
var collection string

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

func initFirebase() {
	// Connect to Firestore
	ctx := context.Background()
	fb := &firebase.Config{}
	app, err := firebase.NewApp(ctx, fb)
	if err != nil {
		log.Fatalln(err)
	}
	fs, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	// TODO provide default config
	collection = os.Getenv("FIREBASE_COLLECTION")
	if collection == "" {
		log.Fatalln("No Firebase collection specified - do so with the environment variable FIREBASE_COLLECTION")
	}
}

func catalogItems(ctx context.Context) (result []Product, err error) {
	iter := fs.Collection(collection).Documents(ctx)
	i := 0
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return result, err
		}
		var p Product
		err = doc.DataTo(&p)
		if err != nil {
			return result, err
		}
		result[i] = p
		i++
	}
	return result, nil
}

func catalogItem(ctx context.Context, sku string) (result Product, err error) {
	doc, err := fs.Collection(collection).Doc(sku).Get(ctx)
	if err != nil {
		return result, err
	}
	err = doc.DataTo(&result)
	return result, err
}

func randomItems(ctx context.Context, count int) (result []Product, err error) {
	result = make([]Product, count)
	for i := 0; i < count; i++ {
		sku := indices[rand.Intn(iLength)]
		item, err := catalogItem(ctx, sku)
		if err != nil {
			fmt.Printf("Error: %v (Sku: %v)", err, sku)
			continue
		}
		result[i] = item
	}
	return result, nil
}
