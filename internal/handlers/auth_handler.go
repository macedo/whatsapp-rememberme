package handlers

import (
	"log"
	"net/http"

	"github.com/macedo/whatsapp-rememberme/internal/app"
)

func SignInPageHandler(c *app.Context) {
	// if c.Session().Exists(c.Context, "user_id") {
	// 	c.Redirect("/admin")
	// 	return
	// }

	if err := c.Render(http.StatusOK, rr.HTML("sign_in")); err != nil {
		log.Println(err)
		c.Response().Write([]byte(err.Error()))
	}
}

// func SignInHandler(c *app.Context) {
// 	var err error

// 	err = r.ParseForm()
// 	if err != nil {
// 		log.Println(err)
// 		w.Write([]byte(err.Error()))
// 		return
// 	}

// 	user, err := repo.GetUserByUsername(r.Context(), r.Form.Get("username"))
// 	if err != nil {
// 		log.Println(err)
// 		c.Session().Flashes()
// 		if err = renderPage(w, r, "sign_in", ""); err != nil {
// 			log.Println(err)
// 			return
// 		}
// 		return
// 	}

// 	encryptor := hash.NewEncryptor(env.Get("SECRET", "pepper"))

// 	if ok := encryptor.Compare(user.EncryptedPassword, r.Form.Get("password")); !ok {
// 		log.Println(err)
// 		session.Put(r.Context(), "flash", "The email or password is incorrect.")
// 		if err := renderPage(w, r, "sign_in", ""); err != nil {
// 			log.Println(err)
// 			return
// 		}
// 		return
// 	}

// 	session.Put(r.Context(), "user_id", user.ID)
// 	session.Put(r.Context(), "flash", "You've been signed in successfully!")

// 	http.Redirect(w, r, "/admin", http.StatusSeeOther)
// }

// func SignOutHandler(c *app.Context) {
// 	_ = session.RenewToken(r.Context())
// 	_ = session.Destroy(r.Context())
// 	_ = session.RenewToken(r.Context())

// 	session.Put(r.Context(), "flash", "You've been logged out successfully!")
// 	http.Redirect(w, r, "/", http.StatusSeeOther)
// }
