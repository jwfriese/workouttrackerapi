package datamodel

type Workout struct {
	Id        int    `json:"id"`
	Timestamp string `json:"timestamp"`
	Lifts     []int  `json:"lifts"`
	Name      string `json:"name"`
}
