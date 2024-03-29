package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route represents a web route.
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes is an array of all the routes that the applicaion uses.
type Routes []Route

// NewRouter creates the routing for us to use.
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(route.HandlerFunc)
	}

	return router
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"QuoteIndex",
		"GET",
		"/quote",
		QuoteIndex,
	},
	Route{
		"QuoteItem",
		"GET",
		"/quote/{quoteId}",
		QuoteByID,
	},
	Route{
		"AddQuote",
		"POST",
		"/quote",
		AddQuote,
	},
	Route{
		"UpdateQuote",
		"PUT",
		"/quote",
		UpdateQuote,
	},
}
