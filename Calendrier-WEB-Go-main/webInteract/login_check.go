package webInteract

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"strconv"
)
// tplLoginForm représente le modèle HTML pour afficher le formulaire de login
var tplLoginForm *template.Template

//compare le mot de passe entré avec le hash du mot de passe en base de données
func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
//Défini les routes pour le login
func InitFormLogin() {
	fmt.Println("appel de InitFormLogin")
	tplLoginForm = template.Must(template.ParseGlob("static/html/*.gohtml"))
	fmt.Println("Loaded Templates:", tplLoginForm.Templates())
	http.HandleFunc("/static/html/form.gohtml", formLoginHandler)
}

// Credentials représente les identifiants de connexion
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
// getToken récupère les identifiants de connexion depuis le fichier credentials.json
func getToken() (Credentials, error) {
	var creds Credentials

	file, err := os.Open("credentials.json")
	if err != nil {
		//tmp
		return creds, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&creds)
	if err != nil {
		//tmp
		return creds, err
	}
	return creds, nil

}
//défini la clé pour la session utilisateur
var store = sessions.NewCookieStore([]byte("t0p-s3cr3t"))

// loginHandler récupère les données du formulaire de login et vérifie si l'utilisateur existe
func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("appel de checkInputFrom")

	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée}", http.StatusMethodNotAllowed)
		fmt.Println(r.URL.Path)
		return
	}

	fmt.Println("Méthode:", r.Method)
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Erreur de parsing du formulaire", http.StatusBadRequest)
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	username = template.HTMLEscapeString(username)
	password = template.HTMLEscapeString(password)

	if username == "" || password == "" {

		redirectURL := "/static/html/form.gohtml?error=champs vides"
		fmt.Printf("Redirection vers: %s\n", redirectURL)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}
	url := "http://dedream.fr/api/user?username=" + username
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	token, err := getToken()
	if err != nil {
		fmt.Println(err)
		return
	}
	req.SetBasicAuth(token.Api.Username, token.Api.Password)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return
	}

	if _, ok := response["error"].(string); ok {
		redirectURL := "/static/html/form.gohtml?error=utilisateur inconnu"
		fmt.Printf("Redirection vers: %s\n", redirectURL)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	if response["username"] != username {
		redirectURL := "/static/html/form.gohtml?error=erreur de login"
		fmt.Printf("Redirection vers: %s\n", redirectURL)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}
	receivedPassword := response["password"].(string)

	if !checkPassword(receivedPassword, password) {
		redirectURL := "/static/html/form.gohtml?error=mot de passe incorrect"
		fmt.Printf("Redirection vers: %s\n", redirectURL)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	println("redirect vers planning.gohtml", w)
	currentDate := time.Now().Format("2006-01-02 15:04:05")
	_, err = time.Parse(time.DateTime, currentDate)
	if err != nil {
		fmt.Println(err)
		return
	}

	id := response["id"].(float64)

	session, err := store.Get(r, "user-session")
	if err != nil {
		fmt.Println(err)
		return
	}
	session.Values["username"] = username
	session.Values["id"] = id

	cookie := http.Cookie{
		Name:  "id",
		Value: strconv.Itoa(int(id)),
		MaxAge:   3600,
	}
	http.SetCookie(w, &cookie)
	session.Save(r, w)

	http.Redirect(w, r, "/planning?date="+currentDate, http.StatusSeeOther)
}
// formLoginHandler affiche le formulaire de login
func formLoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("appel de formLoginHandler")
	err := tplRegisterForm.ExecuteTemplate(w, "form.gohtml", map[string]interface{}{
		"Error":   r.URL.Query().Get("error"),
		"Success": r.URL.Query().Get("success"),
	})

	if err != nil {
		log.Fatalln("Problème de template", err)
	}
}
