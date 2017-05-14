// usage: go run main.go

package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/bolt"
	"github.com/cayleygraph/cayley/quad"
	"github.com/cayleygraph/cayley/schema"
)

var dbPath = "/tmp/db.boltdb"

type Admin struct {
	ID             quad.IRI `quad:"@id"`
	Name           string   `json:"name" quad:"name"`
	Email          string   `json:"email" quad:"email"`
	HashedPassword string   `json:"hashedPassword"  quad:"hashed_password"`
}

func init() {
	schema.RegisterType("Admin", Admin{})
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func genID() quad.IRI {
	return quad.IRI(fmt.Sprintf("%x", rand.Intn(0xffff)))
}

func main() {
	store := initializeAndOpenGraph(dbPath)

	printAllQuads(store)
	printAllAdmins(store)
}

// helper functions

func printAllQuads(store *cayley.Handle) {
	it := store.QuadsAllIterator()
	defer it.Close()
	fmt.Println("\nquads:")
	for it.Next() {
		fmt.Println(store.Quad(it.Result()))
	}
	fmt.Println()
}

func printAllAdmins(store *cayley.Handle) {
	// get all admins
	var admins []Admin
	checkErr(schema.LoadTo(nil, store, &admins))
	fmt.Println("admins:")
	for _, a := range admins {
		fmt.Printf("%+v\n", a)
	}
	fmt.Println()
}

func initializeAndOpenGraph(dbFile string) *cayley.Handle {
	graph.InitQuadStore("bolt", dbFile, nil)

	// Open and use the database
	store, err := cayley.NewGraph("bolt", dbFile, nil)
	if err != nil {
		log.Fatalln(err)
	}

	return store
}
