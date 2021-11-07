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

func home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "welcome to cropper!")
	}
}

func newLink(s service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// obtain request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			err = fmt.Errorf("unable to parse request body %w", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			// slogger.Debug(err)
			// JSONError(rw, err, http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// compose dataset
		link := new(cropper.Link)
		err = json.Unmarshal(body, link)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal JSON %w", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			// slogger.Debug(err)
			// JSONError(rw, err, http.StatusBadRequest)
			return
		}

		// call service level
		err = s.NewLink(link)
		if err != nil {
			err = fmt.Errorf("unable to create link pair: %w", err)
			http.Error(w, "unable to create link pair", http.StatusInsufficientStorage)
			// slogger.Debug(err)
			// JSONError(rw, err, http.StatusBadRequest)
			return
		}

		// slogger.Infof("a new short ID added to database %+v", link)
		// output data
		w.Header().Set("Application", "Cropper")
		w.WriteHeader(http.StatusOK)
	}
}

func redirect(s service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// define shortID from users query
		params := mux.Vars(r)
		short := params["short"]
		// defint corresponding long URL from database
		long, err := s.Resolve(short)
		if err != nil {
			err = fmt.Errorf("cannot obtain ling URI: %w", err)
			http.Error(w, "cannot obtain ling URI", http.StatusInsufficientStorage)
			// slogger.Debugf("resolving error %w", err)
			// JSONError(rw, err, http.StatusBadRequest)
			return
		}
		// slogger.Debugf("successful redirect %s -> %s", shortID, longURL)

		// implement redirect
		http.Redirect(w, r, long, http.StatusPermanentRedirect)
	}
}

func deleteLink(s service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// obtain request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			err = fmt.Errorf("unable to parse request body %w", err)
			http.Error(w, "unable to parse request body", http.StatusBadRequest)
			// slogger.Debug(err)
			// JSONError(rw, err, http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// compose dataset
		link := new(cropper.Link)
		err = json.Unmarshal(body, link)
		if err != nil {
			err = fmt.Errorf("unable to unmarshal JSON %w", err)
			http.Error(w, "unable to unmarshal JSON", http.StatusBadRequest)
			// slogger.Debug(err)
			// JSONError(rw, err, http.StatusBadRequest)
			return
		}

		if link.Short == "" {
			err = fmt.Errorf("the short link is empty, nothing to delete")
			http.Error(w, "the short link is empty, nothing to delete", http.StatusBadRequest)
			// slogger.Debug(err)
			// JSONError(rw, err, http.StatusBadRequest)
			return
		}

		// TODO - execute same check for the user (mean check authorized user deletes his own link)

		// call service level
		userID := uuid.New()
		err = s.DeleteLink(link.Short, &userID)
		if err != nil {
			err = fmt.Errorf("unable to create link pair: %w", err)
			http.Error(w, "unable to create link pair", http.StatusInsufficientStorage)
			// slogger.Debug(err)
			// JSONError(rw, err, http.StatusBadRequest)
			return
		}

		// slogger.Infof("a new short ID added to database %+v", link)
		// output data
		w.Header().Set("Application", "Cropper")
		w.WriteHeader(http.StatusOK)
	}
}
