package Entity


type Player struct {
	Id int `json:"id"`
	Name string `json:"name"`
	PlayingRole string `json:"playingRole"`
	Team Team `json:"team"`
}
