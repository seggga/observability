package endpoint

import (
	"github.com/gorilla/mux"
)

func RegisterPublicHTTP(s service) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", home()).Methods("GET")
	r.HandleFunc("/{short}", redirect(s)).Methods("GET")
	r.HandleFunc("/new-link", authMiddleware(newLink(s))).Methods("POST")
	r.HandleFunc("/delete-link", authMiddleware(deleteLink(s))).Methods("POST")

	return r
}
