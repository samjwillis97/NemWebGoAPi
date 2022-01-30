package middlewares

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// LoggingMW logs resuts made to the server
// Address Request, Method, URL
func LoggingMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("%s \t %s \t %s", r.RemoteAddr, r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func TimingMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		log.Traceln("Execution Time:", time.Since(startTime))
	})
}
