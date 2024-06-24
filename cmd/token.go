package cmd

import (
	"encoding/json"
	"os"

	"golang.org/x/oauth2"
)

const tokenFile = "token.json"

// saveToken saves the token to a file
func saveToken(token *oauth2.Token) error {
	file, err := os.Create(tokenFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(token)
}

// loadToken loads the token from a file
func loadToken() (*oauth2.Token, error) {
	file, err := os.Open(tokenFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var token oauth2.Token
	err = json.NewDecoder(file).Decode(&token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}
