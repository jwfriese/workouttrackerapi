package datamodel

type Set struct {
	Id            int      `json:"id"`
	DataTemplate  string   `json:"dataTemplate"`
	Lift          int      `json:"lift"`
	Weight        *float32 `json:"weight,omitempty"`
	Height        *float32 `json:"height,omitempty"`
	TimeInSeconds *float32 `json:"timeInSeconds,omitempty"`
	Reps          *int     `json:"reps,omitempty"`
}
