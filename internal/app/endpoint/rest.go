package endpoint

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type RESTService struct {
	s service
	l *zap.Logger
	r *mux.Router
}

func RegisterPublicHTTP(s service, l *zap.Logger) *RESTService {
	rs := &RESTService{
		s: s,
		l: l,
	}
	r := mux.NewRouter()

	getR := r.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/", rs.home)
	getR.HandleFunc("/{short}", rs.redirect)

	postR := r.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/new-link", rs.newLink)
	postR.HandleFunc("/delete-link", rs.deleteLink)
	postR.Use(rs.authMiddleware)

	rs.r = r
	return rs
}
