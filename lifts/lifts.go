package lifts

import (
	"net/http"
)

func EndpointHandler() http.Handler {
	return &lift{}
}

type lift struct{}

func (l *lift) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("All the lifts"))
}
