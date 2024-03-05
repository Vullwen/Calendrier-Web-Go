package webInteract

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"fmt"
	"strings"
	//"strconv"
)

// tplCreateEvent représente le modèle HTML pour afficher la page de création d'événement
var tplCreateEvent *template.Template

//Défini les routes pour la création d'événement
func InitCreateEvent() {
	tplCreateEvent = template.Must(template.ParseGlob("static/html/*.gohtml"))
	http.HandleFunc("/createEvent", renderCreateEvent)
	http.HandleFunc("/createEventData",getValueFormCreate)
}

// renderCreateEvent affiche la page de création d'événement
func renderCreateEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("appel de renderCreateEvent")
	session, err := store.Get(r, "user-session")
	username := session.Values["username"].(string)
	if err != nil {
		fmt.Println("Erreur session:", err)	
		http.Redirect(w, r, "/", http.StatusFound)
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
		Username        string
		Categories []Category
	}{
		Username: 	  username,
		Categories: categories,
	}
	err = tplCreateEvent.ExecuteTemplate(w, "createEvent.gohtml", data)
	if err != nil {
		log.Fatalln("Problème de template", err)
	}
}
//fonction qui récupère les données du formulaire de création d'événement
func getValueFormCreate(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	idUser := int(session.Values["id"].(float64))
	if err != nil {
		fmt.Println("Erreur session:", err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
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
	method := "POST"
	payload := strings.NewReader(fmt.Sprintf(`{
		"title": "%s",
		"description": "%s",
		"localisation": "%s",
		"start_date": "%s",
		"end_date": "%s",
		"event_type": "%s"
	}`, title, description, localisation, startDateTime, endDateTime, eventCategoryStr))
	fmt.Println(title,description,localisation,startDateTime,endDateTime,eventCategoryStr)

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

	dataRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Erreur requête body", err)
		http.Error(w, "Erreur requête body", http.StatusInternalServerError)
		return
	}
	idEventCreate := dataRes[0]
	fmt.Println("idEventCreate",idEventCreate)

	//on lie l'event à l'utilisateur
	url = "http://dedream.fr/api/link"
  	method = "POST"
	payload = strings.NewReader(fmt.Sprintf(`{
		"user_id": "%d",
		"event_id": "%d"
	}`, idUser, idEventCreate))

	client = &http.Client{}
	req, err = http.NewRequest(method, url, payload)
	if err != nil {
		log.Println("Erreur requête:", err)
		http.Error(w, "Erreur requête", http.StatusInternalServerError)
		return
	}
	req.SetBasicAuth(token.Api.Username, token.Api.Password)
	res, err = client.Do(req)
	if err != nil {
		log.Println("Erreur requête", err)
		http.Error(w, "Erreur requête", http.StatusInternalServerError)
		return
	}
	fmt.Println(res,"-------",err)
	http.Redirect(w, r, "/planning", http.StatusFound)
}

