package main

// Importing necessary packages and dependencies.
import (
	"api"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"webInteract"

	"github.com/gorilla/mux"
)

// main function - The entry point of the program.
func main() {
	// Serving static files.
	fs := http.FileServer(http.Dir("static"))

	// Dynamic route handling.
	http.HandleFunc("/", redirectFormHandler) // Default route handler.

	// Handling HTML file requests.
	http.HandleFunc("/{path:.*\\.html}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		path := vars["path"]
		http.ServeFile(w, r, path)
	})

	// API setup
	// Routing for API calls, using the LoadAPI function from the 'api' package.
	http.HandleFunc("/api/", api.LoadAPI)

	// Server setup and initialization.
	webInteract.SetupRoads()
	webInteract.InitFormLogin()
	webInteract.InitFormRegister()
	webInteract.InitPlanning()
	webInteract.InitSettings()
	webInteract.InitCreateEvent()
	webInteract.InitCategorie()
	// Static file handling.
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/static/html", redirectFormHandler)
	http.Handle("/static/html/", http.StripPrefix("/static/html/form.gohtml", fs))

	// Starting HTTP server in a new goroutine.
	fmt.Println("Server started on port 80 (HTTP)")
	if err := http.ListenAndServe(":80", nil); err != nil {
		fmt.Println(err)
	}
}

// redirectFormHandler - Function to redirect requests to the form HTML page.
func redirectFormHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/static/html/form.gohtml", http.StatusSeeOther)
}
