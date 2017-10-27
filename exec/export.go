package main

import (
	"fmt"
	"runtime"

	"log"

	"github.com/dernise/venise/cache"
	"github.com/dernise/venise/database"
	"github.com/dernise/venise/writer"
	"gopkg.in/mgo.v2"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println(" ~===== Venise parser version 0.1 Alpha =====~ ")

	//go parser.StartDetailsRoutine()

	tags := make(map[string][]string)
	tags["amenity"] = []string{
		"bicycle_rental",
		"pharmacy",
		"place_of_worship",
		"stripclub",
	}

	//setup mgo
	mongoDatabaseUrl, err := mgo.Dial("localhost:27017") // TODO parse the address
	if err != nil {
		log.Fatal(err.Error())
	}

	mongoDatabase := mongoDatabaseUrl.DB("venise")
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := database.New(mongoDatabase, tags["amenity"])
	if err != nil {
		log.Fatal(err.Error())
	}

	// Setup cache
	cache := cache.NewOSMCache("bin")
	cache.Open()

	nodeWriter := writer.NewNodesWriter(cache, db)
	nodeWriter.WriteNodes(tags)

	fmt.Println("Nodes done")

	waysWriter := writer.NewWaysWriter(cache, db)
	waysWriter.WriteWays(tags)

	fmt.Println("Ways done")

	relationsWriter := writer.NewRelationsWriter(cache, db)
	relationsWriter.WriteRelations(tags)

	fmt.Println("Relations done")
}
