package main

import (
	"backuper/internal/db"
	"backuper/internal/source"
	"fmt"
	"log"
)

func main() {

	//Connect to db
	database, err := db.NewSqlite("./test.db")
	if err != nil {
		log.Fatalf("failed to start db: %v", err)
	}

	sourceHandler, err := source.NewHandler(database)

	if err != nil {
		panic(err)
	}

	err = sourceHandler.CreateFromDsn("postgres://postgres:123123@localhost:5432/postgres")
	if err != nil {
		panic(err)
	}
	//db.Create()
	//Iterate to all connectins get from config or db
	// In the i
	fmt.Println("Hello World")

}
