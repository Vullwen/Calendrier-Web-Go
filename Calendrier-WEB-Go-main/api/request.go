package api

// Importing necessary packages and dependencies.
import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"fmt"
	"strconv"
	
	"golang.org/x/crypto/bcrypt"
)

// Global variable for the database connection.
var db *sql.DB

// Credentials struct to hold database and API credentials.
type Credentials struct {
	Database struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"database"`
	Api struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"api"`
}

type UserEvent struct {
	User int `json:"user"`
	Event int `json:"event"`
}

// Global variable to store loaded credentials.
var creds Credentials

// LoadAPI - Function that handles API requests.
func LoadAPI(w http.ResponseWriter, r *http.Request) {
	// Load credentials.
	creds = loadCredentials(w)

	// Validate credentials.
	if !credentialsAreValid(w, r, creds) {
		// Set header for basic authentication prompt and send unauthorized response.
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted", charset="UTF-8"`)
		errorHandler(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Route the request based on the API path.
	if strings.HasPrefix(r.URL.Path[len("/api/"):], "user") {
		user(w, r)
	} else if strings.HasPrefix(r.URL.Path[len("/api/"):], "event") {
		event(w, r)
	} else if strings.HasPrefix(r.URL.Path[len("/api/"):], "category") {
		category(w, r)
	} else if strings.HasPrefix(r.URL.Path[len("/api/"):], "link") {
		link(w, r)
	} else if strings.HasPrefix(r.URL.Path[len("/api/"):], "dumpDB") {
		dumpDB(w, r)
	} else {
		// Handle invalid API path.
		errorHandler(w, "API not found", http.StatusNotFound)
	}
}

func dumpDB(w http.ResponseWriter, r *http.Request) {
	listTables := []string{"USER", "EVENT", "USER_EVENT", "EVENT_TYPE"}

	// the json struct to hold the data
	type DBDump struct {
		Users []User `json:"users"`
		Events []Event `json:"events"`
		UserEvents []UserEvent `json:"user_events"`
		Category []Category `json:"category"`
	}
	
	var users []User
	var events []Event
	var userEvents []UserEvent
	var categories []Category

	for _, table := range listTables {
		result, err := executeSQLQuery("SELECT * FROM `"+table+"`")
		if err != nil {
			errorHandler(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer result.Close()

		switch table {
		case "USER":
			for result.Next() {

				var user User
				err := result.Scan(&user.ID, &user.Username, &user.Password)
				if err != nil {
					errorHandler(w, "Internal server error : User table", http.StatusInternalServerError)
					return
				}

				users = append(users, user)
			}
		case "EVENT":
			for result.Next() {
				
				var event Event
				err := result.Scan(&event.ID, &event.Title, &event.Description, &event.Localisation, &event.Start_date, &event.End_date, &event.Event_type)
				if err != nil {
					errorHandler(w, "Internal server error : Event table", http.StatusInternalServerError)
					return
				}

				events = append(events, event)
			}

		case "USER_EVENT":
			for result.Next() {
				
				var userEvent UserEvent
				err := result.Scan(&userEvent.User, &userEvent.Event)
				if err != nil {
					errorHandler(w, "Internal server error : UserEvent table", http.StatusInternalServerError)
					return
				}

				userEvents = append(userEvents, userEvent)
			}

		case "EVENT_TYPE":
			for result.Next() {
				
				var category Category
				err := result.Scan(&category.ID, &category.Name)
				if err != nil {
					errorHandler(w, "Internal server error : Category table", http.StatusInternalServerError)
					return
				}

				categories = append(categories, category)
			}

		default:
			errorHandler(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// Marshal the data to JSON and send it as response.
	json.NewEncoder(w).Encode(DBDump{Users: users, Events: events, UserEvents: userEvents, Category: categories})
}

func getIntValue(data map[string]interface{}, key string) int {
    value, _ := strconv.Atoi(fmt.Sprintf("%v", data[key]))
    return value
}

func getStringValue(data map[string]interface{}, key string) string {
    return fmt.Sprintf("%v", data[key])
}

// errorHandler - Function to handle errors and send JSON formatted error messages.
func errorHandler(w http.ResponseWriter, err string, code int) {
	// Set content type and security headers.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	// Set the HTTP status code and encode the error message to JSON.
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{Error: err})
}

func isDateAfter(date1, date2 string) bool {
	t1, err := time.Parse(time.DateTime, date1)
	if err != nil {
		return false
	}

	t2, err := time.Parse(time.DateTime, date2)
	if err != nil {
		return false
	}

	return t1.After(t2)
}

// executeSQLQuery - Function to execute an SQL query.
func executeSQLQuery(query string, args ...interface{}) (*sql.Rows, error) {
	// Open database connection using credentials.
	db, err := sql.Open("mysql", creds.Database.Username+":"+creds.Database.Password+"@tcp(localhost:3306)/DATA")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Prepare the SQL statement.
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the SQL statement with the provided arguments.
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// elementExistsInTable - Function to check if an element exists in a database table.
func elementExistsInTable(table string, attribute string, value string) bool {
	// Execute the query to count occurrences of the value in the specified table and attribute.
	result, err := executeSQLQuery("SELECT COUNT(ID) FROM `"+table+"` WHERE "+attribute+"=?", value)
	if err != nil {
		return false
	}
	defer result.Close()

	// Check if the count is greater than zero.
	if result.Next() {
		var count int
		err := result.Scan(&count)
		if err != nil {
			return false
		}
		return count > 0
	}

	return false
}

// hashPassword - Function to hash a password using bcrypt.
func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// appendLikeCondition appends a 'LIKE' condition to the SQL query if the provided value is non-empty.
// It updates both the SQL query and the parameters slice accordingly.
func appendLikeCondition(request *string, params *[]interface{}, field, value string) {
    if value != "" {
        // Append the 'LIKE' condition to the SQL query and add the parameter to the slice.
        *request += " " + field + " LIKE ? OR"
        *params = append(*params, "%"+value+"%")
    }
}

// appendCondition appends a condition to the SQL query if the provided value is non-empty or if it's a non-strict search.
// It updates both the SQL query and the parameters slice accordingly.
func appendCondition(request *string, params *[]interface{}, field, value string, strictSearch bool) {
    if value != "" {
        // Determine the logical operator based on strict or non-strict search.
        operator := "OR"
        if strictSearch {
            operator = "AND"
        }

        // Additional condition for 'start_date' and 'end_date'.
        if field == "start_date" { 
            *request += " " + field + "<=? " + operator
        } else if field == "end_date" {
            *request += " " + field + ">=? " + operator
        } else {
            *request += " " + field + "=? " + operator
        }

        // Append the parameter to the slice.
        *params = append(*params, value)
    }
}