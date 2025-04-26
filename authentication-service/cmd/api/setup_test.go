package main

import (
	"authentication/data"
	"os"
	"testing"
)


var testApp Config

func TestMain(m *testing.M) {
	// Setup code here
	// This is where you would set up any necessary test environment
	// For example, you might want to set up a test database or mock services

	repo := data.NewPostgresTestRepository(nil)

	testApp.Repo = repo

	// Run the tests
	code := m.Run()

	// Teardown code here
	// This is where you would clean up any resources used during the tests

	os.Exit(code)
}	