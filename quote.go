package main

// Quotes struct which contains an array of quotes. Mind. Blown.
type Quotes struct {
	Quotes []Quote `json:"quotes"`
}

// Quote is a goddammed struct designed to represent a Quote object.
type Quote struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	Author string `json:"author"`
}

// DynamoQuote is a goddammed struct designed to represent a Quote object that comes from AWS DynamoDb.
type DynamoQuote struct {
	QuoteId string `json:"QuoteId"`
	Text    string `json:"text"`
	Author  string `json:"author"`
}
