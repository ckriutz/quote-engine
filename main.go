package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Casey's Quote Engine. Delivering the best in quotes!")

	// Lets make sure we have a table in Dynamo ready to read.
	checkDynamoForQuotesTable()

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))

}
