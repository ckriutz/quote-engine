package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Casey's Quote Engine. Delivering the best in quotes!")

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))

}
