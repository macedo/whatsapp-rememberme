package app

// Auth check user authetication
// func Auth(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if !session.Exists(r.Context(), "user_id") {
// 			http.Redirect(w, r, "/", http.StatusFound)
// 			return
// 		}

// 		w.Header().Set("Cache-Control", "no-store")
// 		next.ServeHTTP(w, r)
// 	})
// }

// // NoSurf implements CSRF protection
// func NoSurf(next http.Handler) http.Handler {
// 	csrfHandler := nosurf.New(next)

// 	csrfHandler.SetBaseCookie(http.Cookie{
// 		HttpOnly: true,
// 		Path:     "/",
// 		Secure:   false,
// 		SameSite: http.SameSiteStrictMode,
// 		Domain:   cfg.GetString("domain"),
// 	})

// 	return csrfHandler
// }

// SessionLoad loads the session on requests
// func SessionLoad(next http.Handler) http.Handler {
// 	return session.LoadAndSave(next)
// }
