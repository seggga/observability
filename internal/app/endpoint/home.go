package endpoint

import (
	"fmt"
	"net/http"
)

func home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "welcome to cropper!")
	}
}
