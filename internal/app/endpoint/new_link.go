package endpoint

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/seggga/observability/pkg/cropper"
)

func newLink(s service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// obtain request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			// err = fmt.Errorf("unable to parse request body %w", err)
			// slogger.Debug(err)
			// JSONError(rw, err, http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// compose dataset
		link := new(cropper.Link)
		err = json.Unmarshal(body, link)
		if err != nil {
			// err = fmt.Errorf("unable to unmarshal JSON %w", err)
			// slogger.Debug(err)
			// JSONError(rw, err, http.StatusBadRequest)
			return
		}

		err = s.NewLink(link)
		if err != nil {
			// err = fmt.Errorf("unable to create link pair: %w", err)
			// slogger.Debug(err)
			// JSONError(rw, err, http.StatusBadRequest)
			return
		}

		// slogger.Infof("a new short ID added to database %+v", link)
		// output data
		rw.Header().Set("Application", "Cropper")
		rw.WriteHeader(http.StatusOK)
	}
}
