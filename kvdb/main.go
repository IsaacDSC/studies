package main

import (
	"fmt"
	"kvdb/pkg"
)

func main() {
	db := pkg.NewKVDB()

	if err := db.CreateCollection("collection1"); err != nil {
		panic(err)
	}

	collection, err := db.GetCollection("collection1")
	if err != nil {
		panic(err)
	}

	if err := collection.CreateFilterKey("filterKey"); err != nil {
		panic(err)
	}

	SetInCollection(collection, "key", map[string]any{
		"data":      "value",
		"filterKey": "filterValue",
	})

	SetInCollection(collection, "key2", map[string]any{
		"data":      "value",
		"filterKey": "filterValue",
	})

	scan, _ := collection.Scan()
	fmt.Println("Scan:", scan)

	// fmt.Println(collection.FindOne(pkg.Key("key")))
	// fmt.Println(collection.FindOne(pkg.Key("key2")))
	// fmt.Println(collection.FindAll(pkg.NewFilterKey("filterKey"), "filterValue"))
}

func SetInCollection(collection *pkg.Collection, key string, input map[string]any) {
	if err := collection.Set(key, input); err != nil {
		panic(err)
	}
}
