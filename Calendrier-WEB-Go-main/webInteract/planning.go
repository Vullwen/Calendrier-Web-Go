package webInteract

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"fmt"
	"time"
)
// structure pour les événements de la semaine, on as pas besoin de toutes les données
type Event struct {
	Id            int
	Title         string
	Description   string
	StartDateTime time.Time
}

// OrganizedEvents représente les événements organisés par jour de la semaine
type OrganizedEvents map[string][]Event

// tplPlanning représente le modèle HTML pour afficher le planning
var tplPlanning *template.Template

// InitPlanning initialise les routes pour le planning et le modèle HTML
func InitPlanning() {
	tplPlanning = template.Must(template.ParseGlob("static/html/*.gohtml"))
	http.HandleFunc("/planning", renderPlanning)
}
// Fonction pour récupérer les événements de l'utilisateur selon l'id de l'utilisateur
func getEventsForUser(id int, r *http.Request) ([]Event, error) {
	// Obtenez le token d'authentification
	token, err := getToken()
	if err != nil {
		return nil, err
	}

	url := "http://dedream.fr/api/link?user_id=" + strconv.Itoa(id)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(token.Api.Username, token.Api.Password)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var rawEvents []struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		StartDate   string `json:"start_date"`
	}

	err = json.Unmarshal(body, &rawEvents)
	if err != nil {
		return nil, err
	}

	var events []Event
	for _, rawEvent := range rawEvents {
		startDateTime, err := time.Parse("2006-01-02 15:04:05", rawEvent.StartDate)
		if err != nil {
			return nil, err
		}

		events = append(events, Event{
			Id:            rawEvent.ID,
			Title:         rawEvent.Title,
			Description:   rawEvent.Description,
			StartDateTime: startDateTime,
		})
	}

	return events, nil
}
// Fonction pour récupérer les événements de la semaine selon l'id de l'utilisateur et la date actuelle
func getEventsForWeek(id int, dateCurrent time.Time, r *http.Request) ([]Event, error) {
	events, err := getEventsForUser(id, r)
	if err != nil {
		return nil, err
	}
	fmt.Println("dateCurrent:", dateCurrent)
	dateStartWeek := dateCurrent
	for dateStartWeek.Weekday() != time.Monday {
		dateStartWeek = dateStartWeek.AddDate(0, 0, -1)
	}


	dateStartWeek = time.Date(dateStartWeek.Year(), dateStartWeek.Month(), dateStartWeek.Day(), 0, 0, 0, 0, dateStartWeek.Location())

	dateEndWeek := dateStartWeek.AddDate(0, 0, 6)
	dateEndWeek = time.Date(dateEndWeek.Year(), dateEndWeek.Month(), dateEndWeek.Day(), 23, 59, 59, 0, dateEndWeek.Location())
	fmt.Println("dateStartWeek:", dateStartWeek)
	fmt.Println("dateEndWeek:", dateEndWeek)
	fmt.Println("events:", events)

	var eventsForWeek []Event
	for _, event := range events {
		if event.StartDateTime.After(dateStartWeek) && event.StartDateTime.Before(dateEndWeek) {
			eventsForWeek = append(eventsForWeek, event)
		}
	}

	return eventsForWeek, nil
}
// Fonction pour traduire le jour de la semaine en français
func translateDay(day string) string {
	switch day {
	case "Monday":
		return "Lundi"
	case "Tuesday":
		return "Mardi"
	case "Wednesday":
		return "Mercredi"
	case "Thursday":
		return "Jeudi"
	case "Friday":
		return "Vendredi"
	case "Saturday":
		return "Samedi"
	case "Sunday":
		return "Dimanche"
	default:
		return day
	}
}
// Fonction pour afficher le planning de l'utilisateur selon la date actuelle
func renderPlanning(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user-session")
	username := session.Values["username"].(string)
	id := int(session.Values["id"].(float64))

	dateParam := r.URL.Query().Get("date")
	var date time.Time
	fmt.Println("dateParam:", dateParam)
	if dateParam != "" {
		date, _ = time.Parse("2006-01-02 15:04:05", dateParam)
	} else {
		date = time.Now().Truncate(24 * time.Hour)
	}
	//fmt.Println("date:", date)
	eventsForWeek, err := getEventsForWeek(id, date, r)

	if err != nil {
		log.Println("Erreur lors de la récupération des événements:", err)
		http.Error(w, "Erreur lors de la récupération des événements", http.StatusInternalServerError)
		return
	}

	organizedEvents := make(OrganizedEvents)
	for _, event := range eventsForWeek {
		eventDate := event.StartDateTime.Truncate(24 * time.Hour)
		day := eventDate.Weekday().String()
		day = translateDay(day)
		//fmt.Println("day:", day)
		if _, ok := organizedEvents[day]; !ok {
			organizedEvents[day] = make([]Event, 0)
		}
		organizedEvents[day] = append(organizedEvents[day], event)
	}

	//fmt.Println(organizedEvents)

	data := struct {
		OrganizedEvents OrganizedEvents
		Username        string
		DaysOfWeek      []string
	}{
		OrganizedEvents: organizedEvents,
		Username:        username,
		DaysOfWeek:      []string{"Lundi", "Mardi", "Mercredi", "Jeudi", "Vendredi", "Samedi", "Dimanche"},
	}

	for _, day := range data.DaysOfWeek {
		log.Printf("Day: %s, Events: %+v", day, data.OrganizedEvents[day])
	}


	err = tplPlanning.ExecuteTemplate(w, "planning.gohtml", data)
	if err != nil {
		log.Fatalln("Problème de template", err)
	}
}


