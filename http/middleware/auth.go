package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/golang/glog"

	"github.com/pruh/api/config"
)

// AuthMiddleware validates basic auth credentials.
func AuthMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc, c *config.Configuration) {
	if c.APIV1Credentials == nil || len(*c.APIV1Credentials) <= 0 {
		glog.Infoln("basic auth users not set, allowing request")
		next(w, r)
		return
	}

	glog.Infoln("checking authentication")
	user, pass, ok := r.BasicAuth()
	if ok && checkCredentials(user, pass, c) {
		glog.Infoln("authentication succeeded")

		next(w, r)
	} else if isLocalNetworkRequest(r, c) {
		glog.Infoln("allow local network request")

		next(w, r)
	} else {
		glog.Infoln("authentication failed")

		w.Header().Set("WWW-Authenticate", `Basic realm="Provide username and password"`)
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("401 Unauthorized.\n"))
		if err != nil {
			glog.Errorf("Cannot write a response. %s", err)
			http.Error(w, "Cannot write a response.", http.StatusInternalServerError)
			return
		}
	}
}

func checkCredentials(user string, password string, c *config.Configuration) bool {
	glog.Infoln("checking credentials")

	if value, ok := (*c.APIV1Credentials)[user]; !ok || value != password {
		glog.Infof("user %s is NOT found or password is incorrect\n", user)
		return false
	}

	glog.Infof("user %s found\n", user)
	return true
}

func isLocalNetworkRequest(r *http.Request, c *config.Configuration) bool {
	remoteIP, err := getRemoteIP(r)
	if err != nil {
		glog.Info(err)
		return false
	}

	headersIP, err := getHeadersIP(r)
	if err != nil {
		glog.Info(err)
		return false
	}

	return isLocalIP(remoteIP, c) && (len(headersIP) == 0 || isLocalIP(headersIP, c))
}

func getRemoteIP(r *http.Request) (net.IP, error) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, fmt.Errorf("Cannot split remote address %s: %v", r.RemoteAddr, err)
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return nil, fmt.Errorf("Cannot parse ip %s", r.RemoteAddr)
	}

	return ip, nil
}

func getHeadersIP(r *http.Request) (net.IP, error) {
	var ip net.IP
	r.Header.Get("")
	if xff := strings.Trim(r.Header.Get("X-Forwarded-For"), ","); len(xff) > 0 {
		addrs := strings.Split(xff, ",")
		lastFwd := strings.TrimSpace(addrs[len(addrs)-1])
		if parsed := net.ParseIP(lastFwd); parsed != nil {
			ip = parsed
		} else {
			return nil, fmt.Errorf("Cannot parse X-Forwarded-For %s", xff)
		}
	} else if xri := r.Header.Get("X-Real-Ip"); len(xri) > 0 {
		if parsed := net.ParseIP(xri); parsed != nil {
			ip = parsed
		} else {
			return nil, fmt.Errorf("Cannot parse X-Real-Ip %s", xri)
		}
	}

	return ip, nil
}

func isLocalIP(ip net.IP, c *config.Configuration) bool {
	glog.Infof("Checking if %s is in local network\n", ip)
	for _, ipnet := range c.LocalNets {
		glog.Infof("Checking network %s\n", ipnet.String())
		if ipnet.Contains(ip) {
			glog.Infof("%s is in local network\n", ip)
			return true
		}
	}

	glog.Infof("%s is NOT in local network\n", ip)
	return false
}
