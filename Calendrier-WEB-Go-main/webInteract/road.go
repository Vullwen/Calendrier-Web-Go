package webInteract

import (
	"fmt"
	"net/http"
)

func SetupRoads() {
	fmt.Println("appel de SetupRoads")
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/logout", logOutSession)
}
