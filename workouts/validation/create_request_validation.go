package validation

import (
	"github.com/jwfriese/workouttrackerapi/workouts/datamodel"
)

type WorkoutsCreateRequestValidator interface {
	Validate([]byte) error
}
