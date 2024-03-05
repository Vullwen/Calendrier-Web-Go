package api

// Importing necessary packages and dependencies.
import (
	"crypto/sha256"
	"encoding/json"
	"net/http"
	"os"
)

// loadCredentials - Function to load API credentials from a JSON file.
func loadCredentials(w http.ResponseWriter) Credentials {
	var creds Credentials // Declare a variable to store credentials.

	// Opening the credentials.json file to read API credentials.
	file, err := os.Open("credentials.json")
	if err != nil {
		// If there's an error opening the file, handle it and return empty credentials.
		errorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		return creds
	}
	defer file.Close() // Ensure file is closed after function execution.

	// Decoding the JSON file into the creds variable.
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&creds)
	if err != nil {
		// Handle JSON decoding errors.
		errorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		return creds
	}

	return creds // Return the loaded credentials.
}

// credentialsAreValid - Function to validate provided credentials against stored ones.
func credentialsAreValid(w http.ResponseWriter, r *http.Request, creds Credentials) bool {
	var usernameHash, passwordHash, expectedUsernameHash, expectedPasswordHash [32]byte

	// Extracting username and password from the HTTP request's Basic Authentication.
	username, password, ok := r.BasicAuth()
	if ok {
		// Hashing the received username and password.
		usernameHash = sha256.Sum256([]byte(username))
		passwordHash = sha256.Sum256([]byte(password))

		// Hashing the expected username and password from stored credentials.
		expectedUsernameHash = sha256.Sum256([]byte(creds.Api.Username))
		expectedPasswordHash = sha256.Sum256([]byte(creds.Api.Password))

		// Comparing the hashes to validate credentials.
		if usernameHash == expectedUsernameHash && passwordHash == expectedPasswordHash {
			return true // Credentials are valid.
		}
	}

	return false // Credentials are invalid.
}
