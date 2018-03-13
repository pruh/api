package middleware

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/j-rooft/api/utils"
)

func AuthMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc, c *utils.Configuration) {
	if c.APIV1Credentials == nil || len(*c.APIV1Credentials) <= 0 {
		log.Println("basic auth users not set, allowing request")
		next(w, r)
		return
	}

	log.Println("checking authentication")
	user, pass, ok := r.BasicAuth()
	if ok && checkCredentials(user, pass, c) {
		log.Println("authentication succeeded")

		next(w, r)
	} else if isLocalNetworkRequest(r, c) {
		log.Println("allow local network request")

		next(w, r)
	} else {
		log.Println("authentication failed")

		w.Header().Set("WWW-Authenticate", `Basic realm="Provide username and password"`)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 Unauthorized.\n"))
	}
}

func checkCredentials(user string, password string, c *utils.Configuration) bool {
	log.Println("checking credentials")

	if value, ok := (*c.APIV1Credentials)[user]; !ok || value != password {
		log.Printf("user %s is NOT found or password is incorrect\n", user)
		return false
	}

	log.Printf("user %s found\n", user)
	return true
}

func isLocalNetworkRequest(r *http.Request, c *utils.Configuration) bool {
	ip, err := getRemoteAddr(r)
	if err != nil {
		return false
	}

	log.Printf("Checking if %s is in local network\n", ip)
	for _, ipnet := range c.LocalNets {
		log.Printf("Checking network %s\n", ipnet.String())
		if ipnet.Contains(ip) {
			log.Printf("%s is in local network\n", ip)
			return true
		}
	}

	log.Printf("%s is NOT in local network\n", ip)
	return false
}

func getRemoteAddr(r *http.Request) (net.IP, error) {
	ip := net.ParseIP(r.RemoteAddr)
	if ip == nil {
		log.Printf("Cannot parse remote address %s", r.RemoteAddr)
		return nil, fmt.Errorf("Cannot parse remote address %s", r.RemoteAddr)
	}

	if xff := strings.Trim(r.Header.Get("X-Forwarded-For"), ","); len(xff) > 0 {
		addrs := strings.Split(xff, ",")
		lastFwd := strings.TrimSpace(addrs[len(addrs)-1])
		println(lastFwd)
		if parsed := net.ParseIP(lastFwd); parsed != nil {
			ip = parsed
		} else {
			log.Printf("Cannot parse X-Forwarded-For %s", xff)
		}
	} else if xri := r.Header.Get("X-Real-Ip"); len(xri) > 0 {
		if parsed := net.ParseIP(xri); parsed != nil {
			ip = parsed
		} else {
			log.Printf("Cannot parse X-Real-Ip %s", xri)
		}
	}

	return ip, nil
}
