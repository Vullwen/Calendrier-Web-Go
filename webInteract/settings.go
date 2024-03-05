package webInteract

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)
// tplSettingsEvent représente le modèle HTML pour afficher les paramètres
var tplSettingsEvent *template.Template

// EventData représente les données d'un événement
type EventData struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	StartDateTime time.Time `json:"start_date"`
	EndDateTime   time.Time `json:"end_date"`
	Location      string    `json:"localisation"`
	TypeEvent     string    `json:"event_type"`
	categoryName  string
}
// Category représente les données d'une catégorie
type Category struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

// getCaetories récupère les catégories depuis l'API
func getCategories() ([]Category, error) {
    token, err := getToken()
    if err != nil {
        return nil, err
    }

    url := "http://dedream.fr/api/category"
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

    var categories []Category
    err = json.Unmarshal(body, &categories)
    if err != nil {
        return nil, err
    }

    return categories, nil
}
//défini les routes pour les paramètres de l'événement
func InitSettings() {
	tplSettingsEvent = template.Must(template.ParseGlob("static/html/*.gohtml"))
	http.HandleFunc("/settings", renderSettings)
	http.HandleFunc("/editEvent", getValueForm)
	http.HandleFunc("/deleteEvent", deleteEventHandler)
}

//getNameCategory récupère le nom de la catégorie depuis l'API en fonction de l'ID
func getNameCategory(id int) string {
	token, err := getToken()
	if err != nil {
		return ""
	}

	url := "http://dedream.fr/api/category?id=" + strconv.Itoa(id)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return ""
	}

	req.SetBasicAuth(token.Api.Username, token.Api.Password)

	res, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ""
	}

	var rawCategory struct {
		Title string `json:"name"`
	}

	err = json.Unmarshal(body, &rawCategory)
	if err != nil {
		return ""
	}

	return rawCategory.Title
}
//getEventById récupère les données d'un événement depuis l'API en fonction de l'ID
func getEventById(id int) (EventData, error) {
	fmt.Println("id", id)
	token, err := getToken()
	if err != nil {
		return EventData{}, err
	}

	url := "http://dedream.fr/api/event?id=" + strconv.Itoa(id)
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return EventData{}, err
	}

	req.SetBasicAuth(token.Api.Username, token.Api.Password)

	res, err := client.Do(req)
	if err != nil {
		return EventData{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return EventData{}, err
	}

	var rawEvents []struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Location    string `json:"localisation"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date"`
		TypeEvent   string  `json:"event_type"`
	}

	fmt.Println("type go of rawEvent", reflect.TypeOf(rawEvents),"\n")
	err = json.Unmarshal(body, &rawEvents)
	if err != nil {
		return EventData{}, err
	}
	fmt.Println("type go of rawEvent", reflect.TypeOf(rawEvents),"\n")
	if err != nil {
		return EventData{}, err
	}

	if len(rawEvents) == 0 {
		return EventData{}, fmt.Errorf("Event not found")
	}

	rawEvent := rawEvents[0]

	startDateTime, err := time.Parse("2006-01-02 15:04:05", rawEvent.StartDate)
	if err != nil {
		return EventData{}, err
	}

	endDateTime, err := time.Parse("2006-01-02 15:04:05", rawEvent.EndDate)
	if err != nil {
		return EventData{}, err
	}

	event := EventData{
		ID:            rawEvent.ID,
		Title:         rawEvent.Title,
		Description:   rawEvent.Description,
		Location:      rawEvent.Location,
		StartDateTime: startDateTime,
		EndDateTime:   endDateTime,
		TypeEvent:     rawEvent.TypeEvent,
	}

	return event, nil
}
//deleteEventHandler lance la suppression d'un événement si on requête la route /deleteEvent
func deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	idEventStr := r.URL.Query().Get("idEvent")
	idEvent, err := strconv.Atoi(idEventStr)
	if err != nil {
		log.Println("Erreur conversion event ID:", err)
		http.Error(w, "Erreur conversion event ID", http.StatusInternalServerError)
		return
	}
	_,err= deleteEventById(idEvent)
	if err != nil {
		log.Println("Erreur lors de la suppression de l'événement:", err)
		http.Error(w, "Erreur lors de la suppression de l'événement", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/planning", http.StatusSeeOther)
	return
}
//deleteEventById supprime un événement depuis l'API en fonction de l'ID
func deleteEventById(id int) (error,error){
	token, err := getToken()
	if err != nil {
		return nil, err
	}

	url := "http://dedream.fr/api/event?id=" + strconv.Itoa(id)
	method := "DELETE"

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
	return nil, nil
}
//renderSettings affiche les données d'un événement en fonction de l'ID
func renderSettings(w http.ResponseWriter, r *http.Request) {
	idEventStr := r.URL.Query().Get("idEvent")
	session, err := store.Get(r, "user-session")
	username := session.Values["username"].(string)
	if err != nil {

		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	idEvent, _ := strconv.Atoi(idEventStr)
	if err != nil {
		log.Println("Erreur lors de la conversion de l'ID de l'événement:", err)
		http.Error(w, "Erreur lors de la conversion de l'ID de l'événement", http.StatusInternalServerError)
		return
	}

	event, err := getEventById(idEvent)
	if err != nil {

		log.Println("Erreur lors de la récupération de l'événement:", err)
		http.Error(w, "Erreur lors de la récupération de l'événement", http.StatusInternalServerError)
		return
	}
	categories, err := getCategories()
    if err != nil {
        log.Println("Erreur des catégories:", err)
        http.Error(w, "Erreur catégories", http.StatusInternalServerError)
        return
    }
	fmt.Println("categories", categories)

	data := struct {
		Event EventData
		Username        string
		Categories []Category
	}{
		Event: event,
		Username: 	  username,
		Categories: categories,
	}
	TypeEventStr,_ := strconv.Atoi(data.Event.TypeEvent)	
	
	data.Event.categoryName = getNameCategory(TypeEventStr)
	fmt.Println(data)

	err = tplSettingsEvent.ExecuteTemplate(w, "settings.gohtml", data)
	if err != nil {
		log.Fatalln("Problème de template", err)
	}
}
//getValueForm récupère les données d'un événement depuis le formulaire si on requête la route /editEvent
func getValueForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Println("Erreur GET form:", err)
		http.Error(w, "Erreur GET form", http.StatusInternalServerError)
		return
	}
	form := r.PostForm
	idEventStr := r.URL.Query().Get("id")
	idEvent, err := strconv.Atoi(idEventStr)
	if err != nil {
		log.Println("Erreur conversion event ID:", err)
		http.Error(w, "Erreur conversion event ID", http.StatusInternalServerError)
		return
	}
	title := form.Get("title")
	description := form.Get("description")
	localisation := form.Get("localisation")
	startDateTime := form.Get("start_date")
	endDateTime := form.Get("end_date")
	eventCategoryStr := form.Get("event_category")

	token, err := getToken()
	if err != nil {
		log.Println("Erreur token:", err)
		http.Error(w, "Erreur token", http.StatusInternalServerError)
		return
	}
	startDateTime = strings.Replace(startDateTime,"T"," ",1)
	endDateTime = strings.Replace(endDateTime,"T"," ",1)

	url := "http://dedream.fr/api/event"
	method := "PUT"
	payload := strings.NewReader(fmt.Sprintf(`{
		"id": %d,
		"title": "%s",
		"description": "%s",
		"localisation": "%s",
		"start_date": "%s",
		"end_date": "%s",
		"event_type": "%s"
	}`, idEvent, title, description, localisation, startDateTime, endDateTime, eventCategoryStr))
	fmt.Println(idEvent,title,description,localisation,startDateTime,endDateTime,eventCategoryStr)

	fmt.Println(url,payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		log.Println("Erreur requête:", err)
		http.Error(w, "Erreur requête", http.StatusInternalServerError)
		return
	}

	req.SetBasicAuth(token.Api.Username, token.Api.Password)

	res, err := client.Do(req)
	if err != nil {
		log.Println("Erreur requête", err)
		http.Error(w, "Erreur requête", http.StatusInternalServerError)
		return
	}
	fmt.Println(res,"-------",err)
	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Erreur requête body", err)
		http.Error(w, "Erreur requête body", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/settings?idEvent="+idEventStr, http.StatusFound)
}

