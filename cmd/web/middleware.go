package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// prevents CSRF attacks to POST requests
// basically says ignore any requests that does not have  CSRF Token
func NoSurfCSRFTokenCheck(next http.Handler) http.Handler {
	csrfCheck := nosurf.New(next)
	csrfCheck.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.UseSecure,
		SameSite: http.SameSiteLaxMode,
		// https://www.youtube.com/watch?v=aUF2QCEudPo
	})
	return csrfCheck
}

// loads and sacves session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
