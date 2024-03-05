package api

// Importing necessary packages and dependencies.
import (
	"database/sql"
	"encoding/json"
	"net/http"
	"fmt"
)

// Category - Struct defining the structure of an category.
type Category struct {
	ID           int    `json:"id"`
	Name		 string `json:"name"`
}

/*
##########################################################################################
#################################### EVENT FUNCTIONS #####################################
##########################################################################################
*/
// category - Main handler for category-related HTTP requests.
func category(w http.ResponseWriter, r *http.Request) {
	// Switch statement to handle different HTTP methods.
	switch r.Method {
	case "GET":
		searchCategory(w, r) // Handle GET requests.
	case "POST":
		createCategory(w, r) // Handle POST requests.
	case "PUT":
		updateCategory(w, r) // Handle PUT requests.
	case "DELETE":
		deleteCategory(w, r) // Handle DELETE requests.
	default:
		errorHandler(w, "Method not allowed", http.StatusMethodNotAllowed) // Handle unsupported methods.
	}
}

// scanCategory scans a row from the result set into an Category struct.
func scanCategory(result *sql.Rows) (Category, error) {
	var category Category
	err := result.Scan(&category.ID, &category.Name)
	return category, err
}

// initCategory - Function to initialize an Category struct from an HTTP request.
func initCategory(r *http.Request) Category {
	// Read request body and unmarshal JSON to a map.
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)

	// Convert map values to strings and populate the Category struct.
	category := Category{
		ID:           getIntValue(data, "id"),
		Name:         getStringValue(data, "name"),
	}
	

	return category
}

func searchCategory(w http.ResponseWriter, r *http.Request) {

	var query string
	var params []interface{}

	// If all the parameters are empty, return all categories.
	if r.URL.Query().Get("id") == "" && r.URL.Query().Get("name") == "" {
		query = "SELECT * FROM EVENT_TYPE"
	} else {
		query = "SELECT * FROM EVENT_TYPE WHERE "
		if r.URL.Query().Get("id") != "" {
			query += "id = ?"
			params = append(params, r.URL.Query().Get("id"))
		}
		if r.URL.Query().Get("name") != "" {
			if r.URL.Query().Get("id") != "" {
				query += " AND "
			}
			query += "name LIKE ?"
			params = append(params, "%"+r.URL.Query().Get("name")+"%")
		}
	}

	result, err := executeSQLQuery(query, params...)
	if err != nil {
		errorHandler(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer result.Close()

	var categories []Category
	for result.Next() {
		category, err := scanCategory(result)
		if err != nil {
			errorHandler(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		categories = append(categories, category)
	}

	w.Header().Set("Content-Type", "application/json ; charset=UTF-8")
	json.NewEncoder(w).Encode(categories)

}

func createCategory(w http.ResponseWriter, r *http.Request) {

	// Initialize an Category struct from the request body.
	category := initCategory(r)

	// Check if the name is empty.
	if category.Name == "" {
		errorHandler(w, "Missing name parameter", http.StatusBadRequest)
		return
	}

	// Check if the category already exists.
	if elementExistsInTable("EVENT_TYPE", "name", category.Name) {
		errorHandler(w, "Category already exists", http.StatusBadRequest)
		return
	}

	// Construct the SQL query to insert the category.
	query := "INSERT INTO `EVENT_TYPE` (`name`) VALUES (?)"
	_, err := executeSQLQuery(query, category.Name)
	if err != nil {
		errorHandler(w, "Failed to execute query", http.StatusInternalServerError)
		return
	}

	// Set response header to JSON.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	message := "Category created successfully"

	// Send the category data as JSON.
	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{Message: message})
}

func updateCategory(w http.ResponseWriter, r *http.Request) {

	// Initialize an Category struct from the request body.
	category := initCategory(r)

	// Check if the name is empty.
	if category.Name == "" || category.ID == 0 {
		errorHandler(w, "Missing name and/or id parameter", http.StatusBadRequest)
		return
	}

	fmt.Println(category)

	// Check if the category already exists.
	if !elementExistsInTable("EVENT_TYPE", "id", string(category.ID)) {
		fmt.Println(category.ID)
		errorHandler(w, "Category does not exist", http.StatusBadRequest)
		return
	}

	// Construct the SQL query to update the category.
	query := "UPDATE `EVENT_TYPE` SET `name` = ? WHERE `id` = ?"
	_, err := executeSQLQuery(query, category.Name, string(category.ID))
	if err != nil {
		errorHandler(w, "Failed to execute query", http.StatusInternalServerError)
		return
	}

	// Set response header to JSON.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	
	message := "Category updated successfully"

	// Send the category data as JSON.
	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{Message: message})
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {

	if r.URL.Query().Get("id") == "" {
		errorHandler(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	// Check if the category already exists.
	if !elementExistsInTable("EVENT_TYPE", "id", string(r.URL.Query().Get("id"))) {
		errorHandler(w, "Category does not exist", http.StatusBadRequest)
		return
	}

	// Construct the SQL query to delete the category.
	query := "DELETE FROM `EVENT_TYPE` WHERE `id` = ?"
	_, err := executeSQLQuery(query, string(r.URL.Query().Get("id")))
	if err != nil {
		errorHandler(w, "Failed to execute query", http.StatusInternalServerError)
		return
	}

	// Set response header to JSON.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	message := "Category deleted successfully"

	// Send the category data as JSON.
	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{Message: message})

}