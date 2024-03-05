package webInteract
import (
	"net/http"
)
//supprime la session de l'utilisateur et le redirige vers la page d'accueil
func logOutSession(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}