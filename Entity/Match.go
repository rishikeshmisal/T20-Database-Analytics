package Entity

type Match struct {
	Id int `json:"id"`
	Details `json:"details"`
	Result `json:"result"`
	Score `json:"score"`
}
