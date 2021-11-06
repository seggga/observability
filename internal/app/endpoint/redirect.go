package endpoint

import (
	"net/http"

	"github.com/gorilla/mux"
)

func redirect(s service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// define shortID from users query
		params := mux.Vars(r)
		short := params["short"]
		// defint corresponding long URL from database
		long, err := s.Resolve(short)
		if err != nil {
			// slogger.Debugf("resolving error %w", err)
			// JSONError(rw, err, http.StatusBadRequest)
			return
		}
		// slogger.Debugf("successful redirect %s -> %s", shortID, longURL)

		// implement redirect
		http.Redirect(rw, r, long, http.StatusPermanentRedirect)
	}
}
