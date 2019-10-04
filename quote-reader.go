package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// readQuoteFromJSONFile does exactly what you think it does.
func readQuotesFromJSONFile() Quotes {
	// First, we open the JSON File
	jsonFile, err := os.Open("quotes.json")

	// if os.Open returns an error, lets do something about it.
	if err != nil {
		fmt.Println("There was an error opening quotes.json. You really screwed this up.")
	}

	// Okay, lets defer this, I guess, so we can parse it here in a bit.
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// Lets initialize the quotes array. For stuff.
	var quotes Quotes

	// We "unmarshal" our byteArray which contains our jsonFile's content into 'quotes' which we defined above.
	json.Unmarshal(byteValue, &quotes)

	// Okay, lets return all of the quotes!
	return quotes
}
