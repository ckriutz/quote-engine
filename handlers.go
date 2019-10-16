package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var quotes Quotes

// Index is the function that is fired when the user goes to the root page.
func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "<h2>Casey's Quote Engine. Delivering the best in quotes!</h2>")
}

// QuoteIndex is the function that fires when the user hits /quote
// At the moment, they would get a random quote.
func QuoteIndex(w http.ResponseWriter, r *http.Request) {
	// First, lets get the quotes from the JSON file. This is in the other file. I *moved it*.
	quotes = readQuotesFromJSONFile()

	// Grab a random number so I can grab a random quote.
	rand.Seed(time.Now().UnixNano())
	var selectedQuoteID = rand.Intn(len(quotes.Quotes))

	// Now, actually get the quote.
	q := &quotes.Quotes[selectedQuoteID]

	// Set all the fun headers.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(q)
}

// QuoteByID is the function that fires when the user hits /quote/{id}
func QuoteByID(w http.ResponseWriter, r *http.Request) {
	// Grab all the querystring things.
	vars := mux.Vars(r)

	// ...specifically the quoteId, which is what the item after the /quote/ is going to be.
	quoteID := vars["quoteId"]

	// Lets convert it into a number, and if that fails, lets just return a 500 error.
	qID, err := strconv.Atoi(quoteID)
	if err != nil {
		// oh! I can return a 500 error!
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Now, lets get the quotes from the JSON file. This is in the other file. I *moved it*.
	quotes = readQuotesFromJSONFile()

	// ...attempt to find the quote the user is looking for.
	q := GetQuoteByID(qID)

	// Set all the fun headers.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(q)
}

// GetQuoteByID will iterate though the quotes, and attempt to find the one with the right ID.
func GetQuoteByID(ID int) Quote {
	for _, qt := range quotes.Quotes {
		if qt.ID == ID {
			return qt
		}
	}
	return Quote{ID: 0, Text: "This did not work.", Author: "Casey"}
}
