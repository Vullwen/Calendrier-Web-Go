package webInteract

import (
	"fmt"
	"html/template"
	"net/http"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	//"reflect"
)
//la structure de la catégorie pour l'affichage
type categorie struct {
	Id   int
	Name string
}
//Défini les routes pour la création de catégorie et sa gestion
func InitCategorie() {
	fmt.Println("appel de InitCategorie")
	tplLoginForm = template.Must(template.ParseGlob("static/html/*.gohtml"))
	fmt.Println("Loaded Templates:", tplLoginForm.Templates())
	http.HandleFunc("/createCategory", categorieHandler)
	http.HandleFunc("/addCategory", addCategoryHandler)
	http.HandleFunc("/deleteCategory", deleteCategoryHandler)
	http.HandleFunc("/updateCategory", editCategoryHandler)
}

// categorieHandler affiche la page de création de catégorie si la bonne url est appelée
func categorieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("appel de categorieHandler")

	// api call to get categories
	query := `http://dedream.fr/api/category`
	method := "GET"
	categories := []categorie{}
	client := &http.Client{}
	req, err := http.NewRequest(method, query, nil)
	if err != nil {
		fmt.Println("Problème de requête", err)
	}

	token, err := getToken()
	if err != nil {
		fmt.Println(err)
		return
	}
	req.SetBasicAuth(token.Api.Username, token.Api.Password)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Problème de requête", err)
	}
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&categories)

	data := struct {
		Categories []categorie
	}{
		Categories: categories,
	}

	err = tplLoginForm.ExecuteTemplate(w, "categorie.gohtml", data)
	if err != nil {
		fmt.Println("Problème de template", err)
	}
}
//handler si l'on souhate supprimer une catégorie
func deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	idEventStr := r.URL.Query().Get("idCat")
	idEvent, err := strconv.Atoi(idEventStr)
	if err != nil {
		log.Println("Erreur conversion event ID:", err)
		http.Error(w, "Erreur conversion event ID", http.StatusInternalServerError)
		return
	}
	_,err= deleteCategory(idEvent)
	if err != nil {
		log.Println("Erreur lors de la suppression de la catégorie:", err)
		http.Error(w, "Erreur lors de la suppression de la catégorie", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/createCategory", http.StatusSeeOther)
	return
}
//supprime une catégorie en fonction de son ID
func deleteCategory(id int) (error,error){
	fmt.Println("appel de deleteCategory")
	token, err := getToken()
	if err != nil {
		return nil, err
	}

	url := "http://dedream.fr/api/category?id=" + strconv.Itoa(id)
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
	fmt.Println(res)
	fmt.Println(err)
	defer res.Body.Close()
	return nil, nil

}

//handler si l'on souhate modifier une catégorie
func editCategoryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("appel de editCategoryHandler")
	idEventStr := r.URL.Query().Get("idCat")
	id,_ := strconv.Atoi(idEventStr)
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
	nameCategory := form.Get("name")
	editCategory(id,nameCategory)
	http.Redirect(w, r, "/createCategory", http.StatusSeeOther)

}
//modifie une catégorie en fonction de son ID et de son nom
func editCategory(id int, name string){
	//api call to edit category
	query := `http://dedream.fr/api/category`
	method := "PUT"
	client := &http.Client{}
	payload := strings.NewReader(fmt.Sprintf(`{
		"name": "%s"
		"id": "%d"
	}`, name, id))
	
	req, err := http.NewRequest(method, query, payload)
	if err != nil {
		fmt.Println("Problème de requête", err)
	}

	token, err := getToken()
	if err != nil {
		fmt.Println(err)
		return
	}
	req.SetBasicAuth(token.Api.Username, token.Api.Password)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Problème de requête", err)
	}
	defer res.Body.Close()
	return
}
//handler si l'on souhate ajouter une catégorie
func addCategoryHandler(w http.ResponseWriter, r *http.Request) {
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
	nameCategory := form.Get("name")
	err:=addCategory(nameCategory)
	if err != nil {
		log.Println("Erreur lors de l'ajout de la catégorie:", err)
		http.Error(w, "Erreur lors de l'ajout de la catégorie", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/createCategory", http.StatusSeeOther)
}
//ajoute une catégorie en fonction de son nom
func addCategory(name string) (error){
	//api call to add category
	fmt.Println("appel de addCategory")

	query := `http://dedream.fr/api/category`
	method := "POST"
	client := &http.Client{}
	payload := strings.NewReader(fmt.Sprintf(`{
		"name": "%s"
	}`, name))

	req, err := http.NewRequest(method, query, payload)
	if err != nil {
		fmt.Println("Problème de requête", err)
	}
	
	token, err := getToken()
	if err != nil {
		fmt.Println(err)
		return err
	} 
	req.SetBasicAuth(token.Api.Username, token.Api.Password)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Problème de requête", err)
	}
	defer res.Body.Close()
	return nil
}