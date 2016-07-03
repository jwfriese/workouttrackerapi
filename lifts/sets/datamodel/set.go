package datamodel

type Set struct {
	Id            int      `json:"-"`
	DataTemplate  string   `json:"-"`
	Lift          int      `json:"-"`
	Weight        *float32 `json:"weight,omitempty"`
	Height        *float32 `json:"height,omitempty"`
	TimeInSeconds *float32 `json:"timeInSeconds,omitempty"`
	Reps          *int     `json:"reps,omitempty"`
}
