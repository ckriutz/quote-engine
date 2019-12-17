package main

import "testing"

func TestConnectToDatabase(t *testing.T) {
	connection, err := connectToDatabase()

	if connection != "success" {
		t.Error("ConnectToDatabase failed.")
	}

	if err != nil {
		t.Error("it fucked up")
	}
}
