package main

import (
	"log"

	"github.com/matias-pan-globant/db"
)

func main() {
	fdb, err := db.NewFileDB("test.data")
	if err != nil {
		panic(err)
	}
	if err := fdb.Create("clave1", "nope"); err == nil {
		log.Fatalf("expected error when creating something with key")
	}
	if err := fdb.Create("clave3", "{\"hellothere\":\"data\"}"); err != nil {
		log.Fatalf("did not expected error: %s", err)
	}
	fdb.Close()
}
