package api

// Importing necessary packages and dependencies.
import (
    "encoding/json"
    "net/http"
    "database/sql"
)

// link - Function to handle link-related API requests.
func link(w http.ResponseWriter, r *http.Request) {
    // Determining the HTTP method and calling the appropriate function.
    switch r.Method {
    case "GET":
        searchLink(w, r) // Handling GET request.
    case "POST":
        createLink(w, r) // Handling POST request.
    case "DELETE":
        deleteLink(w, r) // Handling DELETE request.
    default:
        // Handling invalid HTTP methods.
        errorHandler(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

// searchLink - Function to search for a link between a user and an event in the database.
func searchLink(w http.ResponseWriter, r *http.Request) {
    userId := r.FormValue("user_id")
    eventId := r.FormValue("event_id")
    var query string
    var result *sql.Rows
    var err error
    err = err

    // Set response header to JSON.
    w.Header().Set("Content-Type", "application/json; charset=utf-8")

    switch {
    case userId != "" && eventId != "":
        // Both user ID and event ID are provided, check if the link exists.
        query := "SELECT EXISTS (SELECT 1 FROM `USER_EVENT` WHERE USER = ? AND EVENT = ?)"
        
        result, err = executeSQLQuery(query, userId, eventId)
        if err != nil {
            errorHandler(w, "Internal server error", http.StatusInternalServerError)
            return
        }
        defer result.Close()

        // Scan the result into a boolean.
        var exists bool
        result.Next()
        err := result.Scan(&exists)
        if err != nil {
            errorHandler(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        // if it doesn't exist, return an error.
        if !exists {
            errorHandler(w, "Link does not exist", http.StatusNotFound)
            return
        }

        // Return the result as a JSON object.
        json.NewEncoder(w).Encode(struct {
            Exists bool `json:"exists"`
        }{Exists: exists})

        
    case userId != "":
        // Only user ID is provided, return all events linked to this user.
        query = "SELECT EVENT.* FROM `EVENT` INNER JOIN `USER_EVENT` ON EVENT.ID = USER_EVENT.EVENT WHERE USER_EVENT.USER = ?"
        result, err = executeSQLQuery(query, userId)

    case eventId != "":
        // Only event ID is provided, return all users linked to this event.
        query = "SELECT USER.* FROM `USER` INNER JOIN `USER_EVENT` ON USER.ID = USER_EVENT.USER WHERE USER_EVENT.EVENT = ?"
        result, err = executeSQLQuery(query, eventId)

    default:
        errorHandler(w, "Missing user_id and/or event_id parameter", http.StatusBadRequest)
        return
    }

    if userId != "" {
        // Requête pour les événements
        var events []Event
        for result.Next() {
            var event Event
            if event, err = scanEvent(result); err != nil {
                errorHandler(w, "Error scanning results", http.StatusInternalServerError)
                return
            }
            events = append(events, event)
        }
        json.NewEncoder(w).Encode(events)
    } else if eventId != "" {
        // Requête pour les utilisateurs
        var users []User
        for result.Next() {
            var user User
            if user, err = scanUser(result); err != nil {
                errorHandler(w, "Error scanning results", http.StatusInternalServerError)
                return
            }
            users = append(users, user)
        }
        json.NewEncoder(w).Encode(users)
    }
}

// createLink - Function to create a link between a user and an event in the database.
func createLink(w http.ResponseWriter, r *http.Request) {

    var err error

    // Decode the request body into the User_Event struct.
    body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)
    var User_Event struct {
        User string `json:"user"`
        Event string `json:"event"`
    }
    User_Event.User = getStringValue(data, "user_id")
    User_Event.Event = getStringValue(data, "event_id")

    if User_Event.User == "<nil>" || User_Event.Event == "<nil>" {
        errorHandler(w, "Missing user and/or event parameter", http.StatusBadRequest)
        return
    }

    // Set response header to JSON.
    w.Header().Set("Content-Type", "application/json; charset=utf-8")

    /* // Check if the user and event exist.
    if !elementExistsInTable(`USER`, userId, `1`) {
        errorHandler(w, "User does not exist", http.StatusBadRequest)
        return
    }
    if !elementExistsInTable(`EVENT`, eventId, `1`) {
        errorHandler(w, "Event does not exist", http.StatusBadRequest)
        return
    } */

    // Check if the link already exists.
    query := "SELECT EXISTS (SELECT 1 FROM `USER_EVENT` WHERE USER = ? AND EVENT = ?)"
    result, err := executeSQLQuery(query, User_Event.User, User_Event.Event)
    if err != nil {
        errorHandler(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    defer result.Close()

    // Scan the result into a boolean.
    var exists bool
    result.Next()
    err = result.Scan(&exists)
    if err != nil {
        errorHandler(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    if exists {
        // The link already exists, return an error.
        errorHandler(w, "Link already exists", http.StatusBadRequest)
        return
    }

    // Create the link.
    query = "INSERT INTO `USER_EVENT` (USER, EVENT) VALUES (?, ?)"
    _, err = executeSQLQuery(query, User_Event.User, User_Event.Event)
    if err != nil {
        errorHandler(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // Return a success message.
    json.NewEncoder(w).Encode(struct {
        Success bool `json:"success"`
    }{Success: true})
}

// deleteLink - Function to delete a link between a user and an event in the database.
func deleteLink(w http.ResponseWriter, r *http.Request) {
    userId := r.FormValue("user_id")
    eventId := r.FormValue("event_id")

    // Set response header to JSON.
    w.Header().Set("Content-Type", "application/json; charset=utf-8")

    // Check if the link exists.
    query := "SELECT EXISTS (SELECT 1 FROM `USER_EVENT` WHERE USER = ? AND EVENT = ?)"
    result, err := executeSQLQuery(query, userId, eventId)
    if err != nil {
        errorHandler(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    defer result.Close()

    // Scan the result into a boolean.
    var exists bool
    result.Next()
    err = result.Scan(&exists)
    if err != nil {
        errorHandler(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    if !exists {
        // The link does not exist, return an error.
        errorHandler(w, "Link does not exist", http.StatusBadRequest)
        return
    }

    // Delete the link.
    query = "DELETE FROM `USER_EVENT` WHERE USER = ? AND EVENT = ?"
    _, err = executeSQLQuery(query, userId, eventId)
    if err != nil {
        errorHandler(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // Return a success message.
    json.NewEncoder(w).Encode(struct {
        Success bool `json:"success"`
    }{Success: true})
}