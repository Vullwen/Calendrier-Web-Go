package webInteract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
)
// tplRegisterForm représente le modèle HTML pour afficher le formulaire d'inscription
var tplRegisterForm *template.Template

// InitFormRegister initialise les routes pour le formulaire d'inscription
func InitFormRegister() {
	tplRegisterForm = template.Must(template.ParseGlob("static/html/*.gohtml"))
	http.HandleFunc("/static/html/register.gohtml", formRegisterHandler)
}
// UserCredentials représente les identifiants de connexion
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// registerHandler gère l'inscription d'un utilisateur
func registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("appel de registerHandler")

	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
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
	passwordConf := r.Form.Get("password_conf")

	if username == "" || password == "" || passwordConf == "" {
		redirectURL := "/static/html/register.gohtml?error=champs vides"
		fmt.Printf("Redirection vers: %s\n", redirectURL)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	} else if password != passwordConf {
		redirectURL := "/static/html/register.gohtml?error=les mots de passe ne correspondent pas"
		fmt.Printf("Redirection vers: %s\n", redirectURL)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	passwordRegex := regexp.MustCompile(`^[0-9a-zA-Z!@#$%^.&*()-_+=]{8,}$`)
	if !passwordRegex.MatchString(password) {
		redirectURL := "/static/html/register.gohtml?error=mot de passe invalide au moins 8 caractères et un chiffre"
		fmt.Printf("Redirection vers: %s\n", redirectURL)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	username = template.HTMLEscapeString(username)
	password = template.HTMLEscapeString(password)
	passwordConf = template.HTMLEscapeString(passwordConf)

	userCredentials := UserCredentials{
		Username: username,
		Password: password,
	}

	// Sérialisation de l'objet en JSON
	jsonData, err := json.Marshal(userCredentials)
	if err != nil {
		fmt.Println(err)
		return
	}
	payload := bytes.NewBuffer(jsonData)
	url := "http://dedream.fr/api/user"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	token, err := getToken()
	if err != nil {
		fmt.Println(err)
		return
	}
	req.SetBasicAuth(token.Api.Username, token.Api.Password)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/static/html/register.gohtml?error=erreur de création de compte", http.StatusSeeOther)
		return
	}
	defer res.Body.Close()
	println("redirect vers login.gohtml", w)
	http.Redirect(w, r, "/static/html/form.gohtml?success=Création de compte effectué", http.StatusSeeOther)
}

// formRegisterHandler gère l'affichage du formulaire d'inscription
func formRegisterHandler(w http.ResponseWriter, r *http.Request) {
	err := tplRegisterForm.ExecuteTemplate(w, "register.gohtml", map[string]interface{}{"Error": r.URL.Query().Get("error")})
	if err != nil {
		log.Fatalln("Problème de template", err)
	}
}
