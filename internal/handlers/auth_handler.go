package handlers

import (
	"log"
	"net/http"
)

func SignInPageHandler(w http.ResponseWriter, r *http.Request) {
	if app.Session.Exists(r.Context(), "user_id") {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	if err := renderPage(w, r, "sign_in", ""); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	err = app.Session.RenewToken(r.Context())
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
		return
	}

	err = r.ParseForm()
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
		return
	}

	user, err := repo.GetUserByUsername(r.Form.Get("username"))
	if err != nil {
		log.Println(err)
		app.Session.Put(r.Context(), "flash", "The username or password is incorrect")
		if err = renderPage(w, r, "sign_in", ""); err != nil {
			log.Println(err)
			return
		}
		return
	}

	if ok := app.Encryptor.Compare(user.EncryptedPassword, r.Form.Get("password")); !ok {
		log.Println(err)
		app.Session.Put(r.Context(), "flash", "The email or password is incorrect.")
		if err := renderPage(w, r, "sign_in", ""); err != nil {
			log.Println(err)
			return
		}
		return
	}

	app.Session.Put(r.Context(), "user_id", user.ID)
	app.Session.Put(r.Context(), "flash", "You've been signed in successfully!")

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}