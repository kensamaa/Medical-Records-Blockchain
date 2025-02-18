package services

import (
	"errors"
	"fmt"

	"github.com/kensamaa/blockchain-medical-records/rest-api/utils"
)

// Login simulates user authentication and returns a JWT token.
func Login(username, password string) (string, error) {
	// TODO: Validate credentials against a user store.
	fmt.Println("Simulated login for:", username)
	if username == "admin" && password == "admin" {
		// Generate a token using our utility function.
		return utils.GenerateToken(username)
	}
	return "", errors.New("invalid credentials")
}
