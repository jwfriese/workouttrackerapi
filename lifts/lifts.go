package lifts

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jwfriese/workouttrackerapi/lifts/repository"
)

func LiftsHandler(r repository.LiftRepository) http.Handler {
	return &liftsHandler{
		repository: r,
	}
}

type liftsHandler struct {
	repository repository.LiftRepository
}

func (h *liftsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lifts := h.repository.All()

	liftsJSON, err := json.Marshal(lifts)
	if err != nil {
		log.Fatal(err)
	}

	var bytes bytes.Buffer
	bytes.Write(liftsJSON)

	json := bytes.Bytes()
	w.Write(json)
}
