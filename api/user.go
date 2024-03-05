package api

// Importing necessary packages and dependencies.
import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// User struct - Defines the structure for user data.
type User struct {
	ID       int  `json:"id"`       // User ID, represented as an int64.
	Username string `json:"username"` // Username, represented as a string.
	Password string `json:"password"` // Password, represented as a string.
}

/*
##########################################################################################
#################################### USER FUNCTIONS ######################################
##########################################################################################
*/

// initUser - Function to initialize a User object from an HTTP request.
func initUser(r *http.Request) User {
	// Read request body and unmarshal JSON to a map.
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)

	// Convert map values to strings and populate the Event struct.
	user := User{
		ID:           getIntValue(data, "id"),
		Username:     getStringValue(data, "username"),
		Password:     getStringValue(data, "password"),
	}
	
	return user
}

// user - Function to handle user-related API requests.
func user(w http.ResponseWriter, r *http.Request) {
	// Determining the HTTP method and calling the appropriate function.
	switch r.Method {
	case "GET":
		searchUser(w, r) // Handling GET request.
	case "POST":
		createUser(w, r) // Handling POST request.
	case "PUT":
		updateUser(w, r) // Handling PUT request.
	case "DELETE":
		deleteUser(w, r) // Handling DELETE request.
	default:
		// Handling unsupported HTTP methods.
		errorHandler(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func scanUser(result *sql.Rows) (User, error) {
	var user User
	err := result.Scan(&user.ID, &user.Username, &user.Password)
	return user, err
}

// searchUser - Function to search for a user in the database.
func searchUser(w http.ResponseWriter, r *http.Request) {
	// Checking if the required parameters (id or username) are provided.
	if r.URL.Query().Get("id") == "" && r.URL.Query().Get("username") == "" && r.URL.Query().Get("password") == "" {
		errorHandler(w, "Wrong OR missing parameter", http.StatusBadRequest)
		return
	}

	// Getting parameters from the request.
	id := r.URL.Query().Get("id")
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	// Checking if the user exists in the database.
	if !elementExistsInTable("USER", "id", id) && !elementExistsInTable("USER", "username", username) {
		errorHandler(w, "User not found", http.StatusNotFound)
		return
	}

	// Executing SQL query to retrieve user data.
	result, err := executeSQLQuery("SELECT * FROM `USER` WHERE id=? OR username=?", id, username)
	if err != nil {
		errorHandler(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Moving to the next result (first row).
	result.Next()

	// Creating a User object and scanning the SQL result into it.
	var user User
	err = result.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		errorHandler(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if password != "" {
		if !checkPasswordHash(password, user.Password) {
			errorHandler(w, "Wrong password", http.StatusBadRequest)
			return
		}
	}

	// Returning the user object as a JSON response.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(user)
}

// createUser - Function to create a new user in the database.
func createUser(w http.ResponseWriter, r *http.Request) {
	// Initializing a User object from the request.
	user := initUser(r)

	// Validating that username and password are provided.
	if user.Username == "" || user.Password == "" || user.ID != 0 {
		errorHandler(w, "Wrong OR missing parameter", http.StatusBadRequest)
		return
	}

	// Checking if the username already exists in the database.
	if elementExistsInTable("USER", "username", user.Username) {
		errorHandler(w, "Username already taken", http.StatusBadRequest)
		return
	}

	// Hashing the password.
	user.Password = hashPassword(user.Password)

	// Executing SQL query to insert the new user.
	result, err := executeSQLQuery("INSERT INTO `USER`(username, password) VALUES(?,?)", user.Username, user.Password)
	if err != nil {
		errorHandler(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieving the ID of the newly created user.
	result, err = executeSQLQuery("SELECT ID FROM `USER` WHERE username=? LIMIT 1", user.Username)
	if err != nil {
		errorHandler(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result.Next()

	// Scanning the ID into the User object.
	err = result.Scan(&user.ID)
	if err != nil {
		errorHandler(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Returning the User object as a JSON response.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(user)
}

// updateUser - Function to update an existing user in the database.
func updateUser(w http.ResponseWriter, r *http.Request) {
	// Initializing a User object from the request.
	user := initUser(r)

	if user.ID == 0 {
		errorHandler(w, "Wrong OR missing parameter", http.StatusBadRequest)
		return
	}

	// Validating that username and password are provided.
	if user.Username == "<nil>" && user.Password == "<nil>" {
		errorHandler(w, "Wrong OR missing parameter", http.StatusBadRequest)
		return
	}

	// Checking if the user exists in the database.
	if !elementExistsInTable("USER", "id", strconv.Itoa(int(user.ID))) {
		errorHandler(w, "User not found", http.StatusNotFound)
		return
	}

	// Checking if the username is already taken by another user.
	if elementExistsInTable("USER", "username", user.Username) && !elementExistsInTable("USER", "id", strconv.Itoa(int(user.ID))) {
		errorHandler(w, "Username already taken", http.StatusBadRequest)
		return
	}

	request := "UPDATE `USER` SET PASSWORD=?, USERNAME=? WHERE id=?"
	params := []interface{}{}

	// Hashing the password if it is provided.
	if user.Password != "<nil>" {
		user.Password = hashPassword(user.Password)
	} else {
		result, err := executeSQLQuery("SELECT PASSWORD FROM `USER` WHERE id=?", user.ID)
		if !result.Next() {
			errorHandler(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = result.Scan(&user.Password)
		if err != nil {
			errorHandler(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	params = append(params, user.Password)

	// If the username is not provided, retrieve the user's information before updating.
	if user.Username == "<nil>" {
		request = strings.TrimSuffix(request, ",")
		result, err := executeSQLQuery("SELECT USERNAME FROM `USER` WHERE id=?", user.ID)
		if !result.Next() {
			errorHandler(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = result.Scan(&user.Username)
		if err != nil {
			errorHandler(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	params = append(params, user.Username)

	// Executing SQL query to update the user's information.
	_, err := executeSQLQuery(request, append(params, user.ID)...)
	if err != nil {
		errorHandler(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Returning the updated user object as a JSON response.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(user)
}

// deleteUser - Function to delete a user from the database.
func deleteUser(w http.ResponseWriter, r *http.Request) {
	// Validating that the user ID is provided.
	if r.URL.Query().Get("id") == "" {
		errorHandler(w, "Wrong OR missing parameter", http.StatusBadRequest)
		return
	}

	// Getting the user ID from the request.
	id := r.URL.Query().Get("id")

	// Checking if the user exists in the database.
	if !elementExistsInTable("USER", "id", id) {
		errorHandler(w, "User not found", http.StatusNotFound)
		return
	}

	// Retrieving the user's information before deletion.
	result, err := executeSQLQuery("SELECT * FROM `USER` WHERE id=?", id)
	if !result.Next() {
		errorHandler(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Creating a User object and scanning the SQL result into it.
	var user User
	err = result.Scan(&user.ID, &user.Username, &user.Password)

	// We first need to delete the user's events.
	_, err = executeSQLQuery("DELETE FROM `USER_EVENT` WHERE user=?", id)
	if err != nil {
		errorHandler(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If the event has no longer any participants, we delete it.
	_, err = executeSQLQuery("DELETE FROM `EVENT` WHERE id NOT IN (SELECT event FROM `USER_EVENT`)")
	if err != nil {
		errorHandler(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Executing SQL query to delete the user.
	result, err = executeSQLQuery("DELETE FROM `USER` WHERE id=?", id)
	if err != nil {
		errorHandler(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Returning the deleted user's information as a JSON response.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(user)
}
