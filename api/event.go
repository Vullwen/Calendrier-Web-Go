package api

// Importing necessary packages and dependencies.
import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// Event - Struct defining the structure of an event.
type Event struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Localisation string `json:"localisation"`
	Start_date   string `json:"start_date"`
	End_date     string `json:"end_date"`
	Event_type   string `json:"event_type"`
}

/*
##########################################################################################
#################################### EVENT FUNCTIONS #####################################
##########################################################################################
*/
// event - Main handler for event-related HTTP requests.
func event(w http.ResponseWriter, r *http.Request) {
	// Switch statement to handle different HTTP methods.
	switch r.Method {
	case "GET":
		searchEvent(w, r) // Handle GET requests.
	case "POST":
		createEvent(w, r) // Handle POST requests.
	case "PUT":
		updateEvent(w, r) // Handle PUT requests.
	case "DELETE":
		deleteEvent(w, r) // Handle DELETE requests.
	default:
		errorHandler(w, "Method not allowed", http.StatusMethodNotAllowed) // Handle unsupported methods.
	}
}

// scanEvent scans a row from the result set into an Event struct.
func scanEvent(result *sql.Rows) (Event, error) {
	var event Event
	err := result.Scan(&event.ID, &event.Title, &event.Description, &event.Localisation, &event.Start_date, &event.End_date, &event.Event_type)
	return event, err
}

// initEvent - Function to initialize an Event struct from an HTTP request.
func initEvent(r *http.Request) Event {
	// Read request body and unmarshal JSON to a map.
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)

	// Convert map values to strings and populate the Event struct.
	event := Event{
		ID:           getIntValue(data, "id"),
		Title:        getStringValue(data, "title"),
		Description:  getStringValue(data, "description"),
		Localisation: getStringValue(data, "localisation"),
		Start_date:   getStringValue(data, "start_date"),
		End_date:     getStringValue(data, "end_date"),
		Event_type:   getStringValue(data, "event_type"),
	}
	

	return event
}

// searchEvent - Function to handle the search for events based on various parameters.
func searchEvent(w http.ResponseWriter, r *http.Request) {
	// First, check if any of the search parameters are provided in the request.
	// If none are provided, return a bad request error.
	if r.FormValue("id") == "" && r.FormValue("title") == "" && r.FormValue("description") == "" &&
		r.FormValue("localisation") == "" && r.FormValue("start_date") == "" && r.FormValue("end_date") == "" &&
		r.FormValue("event_type") == "" {
		errorHandler(w, "Wrong OR missing parameter", http.StatusBadRequest)
		return
	}

	// Extract search parameters from the request.
	id := r.FormValue("id")
	title := r.FormValue("title")
	description := r.FormValue("description")
	localisation := r.FormValue("localisation")
	start_date := r.FormValue("start_date")
	end_date := r.FormValue("end_date")
	event_type := r.FormValue("event_type")

	// Initialize variables for constructing the SQL query and parameters.
	var request string
	var params []interface{}
	var result *sql.Rows
	var err error

	if isDateAfter(start_date, end_date) {
		errorHandler(w, "End date is before start date", http.StatusBadRequest)
		return
	}

	// Determine the type of search (strict or large) based on the request parameters.
	// Construct the SQL query accordingly.
	strictSearch := r.FormValue("strictness") == "strict"
	request = "SELECT * FROM `EVENT` WHERE"

	if r.FormValue("like") == "true" {
		// For searches using the 'LIKE' keyword for partial matching.
		appendLikeCondition(&request, &params, "ID", id)
		appendLikeCondition(&request, &params, "title", title)
		appendLikeCondition(&request, &params, "description", description)
		appendLikeCondition(&request, &params, "localisation", localisation)
	} else {
		// For strict search, match exact values or allow empty values for flexibility.
		appendCondition(&request, &params, "ID", id, strictSearch)
		appendCondition(&request, &params, "title", title, strictSearch)
		appendCondition(&request, &params, "description", description, strictSearch)
		appendCondition(&request, &params, "localisation", localisation, strictSearch)
		appendCondition(&request, &params, "start_date", start_date, strictSearch)
		appendCondition(&request, &params, "end_date", end_date, strictSearch)
		appendCondition(&request, &params, "event_type", event_type, strictSearch)
	}

	// Trim the trailing " AND" or " OR" from the query.
	request = strings.TrimSuffix(request, " AND")
	request = strings.TrimSuffix(request, " OR")

	// Execute the constructed SQL query.
	result, err = executeSQLQuery(request, params...)
	if err != nil {
		// Handle any errors that occur during query execution.
		errorHandler(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a slice to hold the results.
	var events []Event
	for result.Next() {
		// Use the scanEvent function to scan the result into the event struct.
		event, err := scanEvent(result)
		if err != nil {
			// Handle any errors that occur while scanning the result.
			errorHandler(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Append the event struct to the events slice.
		events = append(events, event)
	}

	// Check if any events were found. If not, return a not found error.
	if len(events) == 0 {
		errorHandler(w, "Event not found", http.StatusNotFound)
		return
	}

	// Set the header to indicate the content type as JSON.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// Encode the events slice into JSON and write it to the response.
	json.NewEncoder(w).Encode(events)
}

// createEvent - Function to create a new event based on the HTTP request.
func createEvent(w http.ResponseWriter, r *http.Request) {
	// Initialize an Event struct from the request.
	var event Event
	event = initEvent(r)

	// Check if all necessary parameters are provided.
	// Return an error if any mandatory parameter is missing.
	if event.Title == "" || event.Description == "" || event.Localisation == "" || event.Start_date == "" || event.End_date == "" || event.Event_type == "" {
		errorHandler(w, "Missing parameter", http.StatusBadRequest)
		return
	}

	// Check if the provided event type exists in the database.
	// If it does not exist, return a not found error.
	if !elementExistsInTable("EVENT_TYPE", "id", event.Event_type) {
		errorHandler(w, "Event type not found", http.StatusNotFound)
		return
	}

	if isDateAfter(event.Start_date, event.End_date) {
		errorHandler(w, "End date is before start date", http.StatusBadRequest)
		return
	}

	// Insert the new event into the database.
	// The SQL query uses parameter placeholders for security (preventing SQL injection).
	_, err := executeSQLQuery("INSERT INTO `EVENT` (title, description, localisation, start_date, end_date, event_type) VALUES (?, ?, ?, STR_TO_DATE(?, '%Y-%m-%d %H:%i:%s'), STR_TO_DATE(?, '%Y-%m-%d %H:%i:%s'), ?)", event.Title, event.Description, event.Localisation, event.Start_date, event.End_date, event.Event_type)
	if err != nil {
		// Handle any SQL execution errors.
		errorHandler(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// we get the id of the event we just created
	result, err := executeSQLQuery("SELECT ID FROM EVENT ORDER BY ID DESC LIMIT ?", 1)
	if err != nil {
		errorHandler(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result.Next()
	err = result.Scan(&event.ID)
	if err != nil {
		errorHandler(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the header to indicate the content type as JSON.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// Encode the created event into JSON and write it to the response.
	json.NewEncoder(w).Encode(event)
}

// updateEvent - Function to update an existing event.
func updateEvent(w http.ResponseWriter, r *http.Request) {
	// Initialize an Event struct from the request.
	var event Event
	event = initEvent(r)

	// Check if the event ID is provided, as it's necessary for updating an event.
	if event.ID == 0 {
		errorHandler(w, "Missing parameter", http.StatusBadRequest)
		return
	}

	// Initialize variables to construct the SQL update query.
	var request string
	var params []interface{}

	// Check each field of the Event struct. If a field is not "<nil>", it means it needs to be updated.
	// Append the corresponding SQL fragment and parameter for each field that needs updating.
	if event.Title != "<nil>" {
		request += " title=?,"
		params = append(params, event.Title)
	}
	if event.Description != "<nil>" {
		request += " description=?,"
		params = append(params, event.Description)
	}
	if event.Localisation != "<nil>" {
		request += " localisation=?,"
		params = append(params, event.Localisation)
	}
	if event.Start_date != "<nil>" {
		request += " start_date=STR_TO_DATE(?, '%Y-%m-%d %H:%i:%s'),"
		params = append(params, event.Start_date)
	}
	if event.End_date != "<nil>" {
		request += " end_date=STR_TO_DATE(?, '%Y-%m-%d %H:%i:%s'),"
		params = append(params, event.End_date)
	}
	if event.Event_type != "<nil>" {
		// Check if the event to be updated exists in the database.
		if !elementExistsInTable("EVENT", "id", strconv.Itoa(event.ID)) {
			errorHandler(w, "Event not found", http.StatusNotFound)
			return
		}
		request += " event_type=?,"
		params = append(params, event.Event_type)
	}

	// Remove the trailing comma from the SQL query string.
	request = request[:len(request)-1]

	// Check for date consistency (start date should be before end date).
	// Retrieve current start and end dates from the database to compare with the provided dates.
	result, err := executeSQLQuery("SELECT start_date, end_date FROM `EVENT` WHERE id=?", event.ID)
	if err != nil {
		errorHandler(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result.Next()
	var oldEvent Event
	err = result.Scan(&oldEvent.Start_date, &oldEvent.End_date)
	if err != nil {
		errorHandler(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if (event.Start_date != "<nil>" && event.End_date != "<nil>" && event.Start_date > event.End_date) || (event.Start_date == "<nil>" && event.End_date != "<nil>" && event.End_date < oldEvent.Start_date) || (event.End_date == "<nil>" && event.Start_date != "<nil>" && event.Start_date > oldEvent.End_date) {
		errorHandler(w, "End date is before start date", http.StatusBadRequest)
		return
	}

	// Execute the SQL query to update the event.
	_, err = executeSQLQuery("UPDATE `EVENT` SET" + request + " WHERE id=?", append(params, event.ID)...)
	if err != nil {
		errorHandler(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the response header to JSON.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// Encode the updated event into JSON and write it to the response.
	json.NewEncoder(w).Encode(event)
}

// deleteEvent - Function to delete an existing event.
func deleteEvent(w http.ResponseWriter, r *http.Request) {
	// Check if the 'id' parameter is provided in the request.
	// This parameter is necessary to identify which event to delete.
	if r.FormValue("id") == "" {
		// If the 'id' parameter is missing, return a bad request error.
		errorHandler(w, "Wrong OR missing parameter", http.StatusBadRequest)
		return
	}

	// Extract the 'id' parameter from the request.
	id := r.FormValue("id")

	// Check if the event with the given ID exists in the database.
	// If the event does not exist, return a not found error.
	if !elementExistsInTable("EVENT", "id", id) {
		errorHandler(w, "Event not found", http.StatusNotFound)
		return
	}

	// Retrieve the event from the database before deleting it.
	// This step is usually done to return the deleted event data in the response.
	result, err := executeSQLQuery("SELECT * FROM `EVENT` WHERE id=?", id)
	if err != nil || !result.Next() {
		// If there's an error in executing the query or the event is not found, return an error.
		errorHandler(w, "Event not found or query failed", http.StatusInternalServerError)
		return
	}

	// Create a Event object and populate its fields.
	var event Event
	err = result.Scan(&event.ID, &event.Title, &event.Description, &event.Localisation, &event.Start_date, &event.End_date, &event.Event_type)

    // We get the start date and end date of the event before updating it
    result, err = executeSQLQuery("SELECT start_date, end_date FROM `EVENT` WHERE id=?", event.ID)
    if err != nil {
        errorHandler(w, err.Error(), http.StatusInternalServerError)
        return
    }
    result.Next()
    // Create a Event object and populate its fields.
    var oldEvent Event
    err = result.Scan(&oldEvent.Start_date, &oldEvent.End_date)
    if err != nil {
        errorHandler(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
	if event.Start_date != "<nil>" && event.End_date != "<nil>" && isDateAfter(event.Start_date, event.End_date) || event.Start_date == "<nil>" && event.End_date != "<nil>" && isDateAfter(event.End_date, oldEvent.Start_date) || event.End_date == "<nil>" && event.Start_date != "<nil>" && isDateAfter(event.Start_date, oldEvent.End_date) {
		errorHandler(w, "End date is before start date", http.StatusBadRequest)
		return
	}

	_, err = executeSQLQuery("DELETE FROM `USER_EVENT` WHERE event=?", id)
	if err != nil {
		errorHandler(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the SQL query to delete the event from the database.
	_, err = executeSQLQuery("DELETE FROM `EVENT` WHERE id=?", id)
	if err != nil {
		// If there's an error in executing the delete query, return an internal server error.
		errorHandler(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the response header to indicate that the content type is JSON.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// Encode the deleted event into JSON and write it to the response.
	json.NewEncoder(w).Encode(event)
}
