package Entity


type Score struct {

	Innings1_wickets int `json:"innings1Wickets"`
	Innings1_overs_batted float64 `json:"innings1OversBatted"`
	Innings1_overs float64 `json:"innings1Overs"`
	Innings1_score int `json:"innings1Score"`
	Innings2_score int `json:"innings2Score"`
	Innings2_wickets int `json:"innings2Wickets"`
	Innings2_overs_batted float64 `json:"innings2OversBatted"`
	Innings2_overs float64 `json:"innings2Overs"`
	Dl_method int `json:"dlMethod"`
	Target int `json:"target"`
}
