package datamodel

type Set struct {
	Id            int      `json:"-"`
	DataTemplate  string   `json:"-"`
	Lift          int      `json:"-"`
	Weight        *float32 `json:"Weight,omitempty"`
	Height        *float32 `json:"Height,omitempty"`
	TimeInSeconds *float32 `json:"TimeInSeconds,omitempty"`
	Reps          *int     `json:"Reps,omitempty"`
}
