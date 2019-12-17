package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	// First, lets get the number of quotes from DynamoDB
	quoteCount := getCountOfQuotes()

	// Grab a random number so I can grab a random quote.
	rand.Seed(time.Now().UnixNano())
	var selectedQuoteID = rand.Intn(quoteCount)

	// Now, actually get the quote.
	q := getQuoteFromDynamo(selectedQuoteID)

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

	// ...attempt to find the quote the user is looking for.
	q := getQuoteFromDynamo(qID)

	// Set all the fun headers.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(q)
}

func AddQuote(w http.ResponseWriter, r *http.Request) {
	// First, lets try to read this thing in.
	body, err := ioutil.ReadAll(r.Body)

	// Well, something happened.
	if err != nil {
		panic(err)
	}

	// Something else equally awful happened.
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	quote := DynamoQuote{}

	// We got something, sure, but it looks like it was not a quote.
	if err := json.Unmarshal(body, &quote); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	// At this point it looks like we are in a better place. Lets add it to the list.
	quoteCount := getCountOfQuotes() // Since our first quote starts with 0, the next quote would be the count.
	quote.QuoteId = strconv.Itoa(quoteCount)

	addQuoteToDynamo(quote)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(quote); err != nil {
		panic(err)
	}

}

func UpdateQuote(w http.ResponseWriter, r *http.Request) {
	// First, lets try to read this thing in.
	body, err := ioutil.ReadAll(r.Body)

	// Well, something happened.
	if err != nil {
		panic(err)
	}

	// Something else equally awful happened.
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	quote := DynamoQuote{}

	// We got something, sure, but it looks like it was not a quote.
	if err := json.Unmarshal(body, &quote); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	// At this point it looks like we are in a better place. Lets update the list!
	updateQuote(quote)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(quote); err != nil {
		panic(err)
	}
}
