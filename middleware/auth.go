package middleware

import (
	"log"
	"net/http"

	"github.com/j-rooft/api/utils"
)

// AuthMiddleware validates basic auth credentials.
func AuthMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc, config *utils.Configuration) {
	if config.APIV1Credentials == nil || len(*config.APIV1Credentials) <= 0 {
		log.Println("basic auth users not set, allowing request")
		next(w, r)
		return
	}

	log.Println("checking authentication")
	user, pass, ok := r.BasicAuth()
	if ok && checkCredentials(user, pass, config) {
		log.Println("authentication succeeded")

		next(w, r)
	} else {
		log.Println("authentication failed")

		w.Header().Set("WWW-Authenticate", `Basic realm="Provide username and password"`)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 Unauthorized.\n"))
	}
}

func checkCredentials(user string, password string, config *utils.Configuration) bool {
	log.Println("checking credentials")

	if value, ok := (*config.APIV1Credentials)[user]; !ok || value != password {
		log.Printf("user %s NOT found or password incorrect\n", user)
		return false
	}

	log.Printf("user %s found\n", user)
	return true
}
