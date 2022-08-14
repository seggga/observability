package endpoint

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/seggga/observability/pkg/cropper"
)

func (rs *RESTService) home(w http.ResponseWriter, r *http.Request) {
	rs.l.Debug("home func called")

	fmt.Fprint(w, "welcome to cropper!")
}

func (rs *RESTService) redirect(w http.ResponseWriter, r *http.Request) {
	rs.l.Debug("redirect func called")

	// define shortID from users query
	params := mux.Vars(r)
	short := params["short"]
	// define corresponding long URL from database
	long, err := rs.s.Resolve(short)
	if err != nil {
		err = fmt.Errorf("cannot obtain ling URI: %w", err)
		rs.l.Error(err.Error())

		http.Error(w, "cannot obtain ling URI", http.StatusInsufficientStorage)
		// JSONError(rw, err, http.StatusBadRequest)
		return
	}

	rs.l.Sugar().Debugf("successful redirect %s -> %s", short, long)

	http.Redirect(w, r, long, http.StatusPermanentRedirect)
}

// authMiddleware checks user's authorization
func (rs *RESTService) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// check authorization
			user, pass, ok := r.BasicAuth()
			if !ok {
				err := fmt.Errorf("unauthorized access")
				rs.l.Error(err.Error())
				http.Error(w, "unauthorized access", http.StatusUnauthorized)
				// JSONError(rw, err, http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r)
		},
	)
}

func (rs *RESTService) newLink(w http.ResponseWriter, r *http.Request) {
	rs.l.Debug("newLink func called")

	// obtain request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("unable to parse request body %w", err)

		rs.l.Error(err.Error())

		http.Error(w, err.Error(), http.StatusBadRequest)
		// JSONError(rw, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// compose dataset
	link := new(cropper.Link)
	err = json.Unmarshal(body, link)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal JSON %w", err)

		rs.l.Error(err.Error())

		http.Error(w, err.Error(), http.StatusBadRequest)
		// JSONError(rw, err, http.StatusBadRequest)
		return
	}

	// call service level
	err = rs.s.NewLink(link)
	if err != nil {
		err = fmt.Errorf("unable to create link pair: %w", err)

		rs.l.Error(err.Error())

		http.Error(w, "unable to create link pair", http.StatusInsufficientStorage)
		// JSONError(rw, err, http.StatusBadRequest)
		return
	}

	rs.l.Sugar().Infof("a new short ID added to database %+v", link)

	// output data
	w.Header().Set("Application", "Cropper")
	w.WriteHeader(http.StatusOK)
}

func (rs *RESTService) deleteLink(w http.ResponseWriter, r *http.Request) {
	rs.l.Debug("deleteLink func called")

	// obtain request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("unable to parse request body %w", err)

		rs.l.Error(err.Error())

		http.Error(w, "unable to parse request body", http.StatusBadRequest)
		// JSONError(rw, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// compose dataset
	link := new(cropper.Link)
	err = json.Unmarshal(body, link)
	if err != nil {
		err = fmt.Errorf("unable to unmarshal JSON %w", err)

		rs.l.Error(err.Error())

		http.Error(w, "unable to unmarshal JSON", http.StatusBadRequest)
		// JSONError(rw, err, http.StatusBadRequest)
		return
	}

	if link.Short == "" {
		err = fmt.Errorf("the short link is empty, nothing to delete")

		rs.l.Error(err.Error())

		http.Error(w, "the short link is empty, nothing to delete", http.StatusBadRequest)
		// JSONError(rw, err, http.StatusBadRequest)
		return
	}

	// TODO - execute same check for the user (mean check authorized user deletes his own link)

	// call service level
	userID := uuid.New()
	err = rs.s.DeleteLink(link.Short, &userID)
	if err != nil {
		err = fmt.Errorf("unable to create link pair: %w", err)

		rs.l.Error(err.Error())

		http.Error(w, "unable to create link pair", http.StatusInsufficientStorage)
		// JSONError(rw, err, http.StatusBadRequest)
		return
	}

	rs.l.Sugar().Infof("a new short ID added to database %+v", link)

	// output data
	w.Header().Set("Application", "Cropper")
	w.WriteHeader(http.StatusOK)
}
