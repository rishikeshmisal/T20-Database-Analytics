package main

import (
	"T20-Database-Analytics/Entity"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"io/ioutil"
)

var match []Entity.Match
var series []Entity.Series

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/gethightotalchase/{venue}", GetHighTotalChase).Methods("GET")
	router.HandleFunc("/getmatch/{matchid}", GetMatchById).Methods("GET")
	router.HandleFunc("/getseries/{seriesid}", GetSeriesById).Methods("GET")
	router.HandleFunc("/getteammatches/{team}", GetTeamMatches).Methods("GET")
	router.HandleFunc("/getmatchesbyvenue/{venue}", GetMatchesByVenue).Methods("GET")
	router.HandleFunc("/gethighestscorebyvenue/{venue}", GetHighestScoreByVenue).Methods("GET")
	router.HandleFunc("/getlowestscorebyvenue/{venue}", GetLowestScoreByVenue).Methods("GET")
	router.HandleFunc("/build", buildDB).Methods("POST")
	router.HandleFunc("/espn/getseriesmatches/{seriesid}", GetEspnSeriesById).Methods("GET")
	http.ListenAndServe(":9000", router)

}

func GetEspnSeriesById(writer http.ResponseWriter, request *http.Request){
	matchVar := mux.Vars(request)
	seriesId, err := strconv.Atoi(matchVar["seriesid"])
	writer.Header().Add("Content-Type", "application/json")
	if err != nil {
		fmt.Printf("Wrong input bhai!, %s", err)

		writer.Write([]byte("Wrong input bhai!"))
		return
	}

	url := "https://hs-consumer-api.espncricinfo.com/v1/pages/series/schedule?lang=en&seriesId="+strconv.Itoa(seriesId)
	fmt.Println(url)
	req, _ := http.NewRequest("GET", url, nil)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}

	responseData, err := ioutil.ReadAll(response.Body)

	var result map[string]interface{}
	json.Unmarshal([]byte(responseData), &result)



	fmt.Printf("%s",string(responseData))

	content := result["content"].(map[string]interface{})

	matches := content["matches"].([]interface{})

	matchIds := make([]int,0)

	for _, match := range matches{
		mm := match.(map[string]interface{})
		id := mm["objectId"].(float64)
		if mm["stage"].(string) == "FINISHED" {
			matchIds = append(matchIds, int(id))
		}
	}

	url = "https://hs-consumer-api.espncricinfo.com/v1/pages/match/smart-scorecard?lang=en&seriesId="+
		strconv.Itoa(seriesId)+"&matchId="+strconv.Itoa(matchIds[0])

	req, _ = http.NewRequest("GET", url, nil)
	response, err = http.DefaultClient.Do(req)
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}

	responseData, err = ioutil.ReadAll(response.Body)

	writer.Write(responseData)


}

func GetHighTotalChase(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "application/json")
	matchVar := mux.Vars(request)
	venue := matchVar["venue"]
	jsonEncode := json.NewEncoder(writer)
	venue = strings.ToUpper(venue)

	retMatch := make(map[string][]Entity.Match)

	venues := make(map[string]bool)
	max := make(map[string]int)
	for _, m := range match {
		if strings.Contains(strings.ToUpper(m.Venue), venue) {
			if _, ok := venues[m.Venue]; !ok {
				venues[m.Venue] = true
				max[m.Venue] = -1
			}
			if m.Winner == m.Innings2 {
				if m.Dl_method == 0 {
					if m.Innings1_score > max[m.Venue] {
						max[m.Venue] = m.Innings1_score
						retMatch[m.Venue] = retMatch[m.Venue][:0]
						retMatch[m.Venue] = append(retMatch[m.Venue], m)
					} else if m.Innings1_score == max[m.Venue] {
						retMatch[m.Venue] = append(retMatch[m.Venue], m)
					}
				} else {
					if m.Target > max[m.Venue] {
						max[m.Venue] = m.Target
						retMatch[m.Venue] = retMatch[m.Venue][:0]
						retMatch[m.Venue] = append(retMatch[m.Venue], m)
					} else if m.Target == max[m.Venue] {
						retMatch[m.Venue] = append(retMatch[m.Venue], m)
					}
				}
			}
		}
	}

	jsonEncode.Encode(retMatch)
	return

}
func GetLowestScoreByVenue(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "application/json")
	matchVar := mux.Vars(request)
	venue := matchVar["venue"]
	jsonEncode := json.NewEncoder(writer)
	venue = strings.ToUpper(venue)

	retMatch := make(map[string][]Entity.Match)
	min := make(map[string]int)
	venues := make(map[string]bool)
	for _, m := range match {
		if strings.Contains(strings.ToUpper(m.Venue), venue) {
			if _, ok := venues[m.Venue]; !ok {
				venues[m.Venue] = true
				min[m.Venue] = 1000
			}
			if min[m.Venue] > m.Innings1_score && m.Innings1_score > 0 {
				if min[m.Venue] > m.Innings2_score && m.Innings2_score < m.Innings1_score {
					min[m.Venue] = m.Innings2_score
				} else if m.Innings1_score > 0 {
					min[m.Venue] = m.Innings1_score
				}
				retMatch[m.Venue] = retMatch[m.Venue][:0]
				retMatch[m.Venue] = append(retMatch[m.Venue], m)
			} else if min[m.Venue] > m.Innings2_score && m.Innings2_score > 0 {
				min[m.Venue] = m.Innings2_score
				retMatch[m.Venue] = retMatch[m.Venue][:0]
				retMatch[m.Venue] = append(retMatch[m.Venue], m)
			} else if min[m.Venue] == m.Innings1_score || min[m.Venue] == m.Innings2_score {
				if m.Innings1_score > 0 && m.Innings2_score > 0 {
					retMatch[m.Venue] = append(retMatch[m.Venue], m)
				}
			}
		}
	}

	jsonEncode.Encode(retMatch)
	return

}
func GetHighestScoreByVenue(writer http.ResponseWriter, request *http.Request) {

	writer.Header().Add("Content-Type", "application/json")
	matchVar := mux.Vars(request)
	venue := matchVar["venue"]
	jsonEncode := json.NewEncoder(writer)
	venue = strings.ToUpper(venue)

	venues := make(map[string]bool)
	retMatch := make(map[string][]Entity.Match)
	max := make(map[string]int)
	for _, m := range match {
		if strings.Contains(strings.ToUpper(m.Venue), venue) {
			if _, ok := venues[m.Venue]; !ok {
				venues[m.Venue] = true
				max[m.Venue] = -1
			}
			if max[m.Venue] < m.Innings1_score {
				if max[m.Venue] < m.Innings2_score && m.Innings2_score > m.Innings1_score {
					max[m.Venue] = m.Innings2_score
				} else {
					max[m.Venue] = m.Innings1_score
				}
				retMatch[m.Venue] = retMatch[m.Venue][:0]
				retMatch[m.Venue] = append(retMatch[m.Venue], m)
			} else if max[m.Venue] < m.Innings2_score {
				max[m.Venue] = m.Innings2_score
				retMatch[m.Venue] = retMatch[m.Venue][:0]
				retMatch[m.Venue] = append(retMatch[m.Venue], m)
			} else if max[m.Venue] == m.Innings1_score || max[m.Venue] == m.Innings2_score {
				retMatch[m.Venue] = append(retMatch[m.Venue], m)
			}
		}
	}

	jsonEncode.Encode(retMatch)
	return

}
func buildDB(writer http.ResponseWriter, request *http.Request) {
	buildSeriesDB()
	buildMatchesDB()
	if len(series) > 0 && len(match) > 0 {
		writer.Write([]byte(fmt.Sprintf("Done! updating db with size being %d %d", len(series), len(match))))
	} else {
		writer.Write([]byte(fmt.Sprintf("Some problem with updating db %d %d", len(series), len(match))))
	}
}
func buildMatchesDB() {
	f, err := os.Open("/Users/rishikeshmisal/Work/src/T20-Database-Analytics/t20_matches.csv")
	reader := csv.NewReader(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		newMatch := Entity.Match{}
		newResult := Entity.Result{}
		newScore := Entity.Score{}
		newDetails := Entity.Details{}
		newMatch.Id, err = strconv.Atoi(strings.TrimSpace(line[0]))
		ser, err := strconv.Atoi(strings.TrimSpace(line[1]))
		newDetails.Date = strings.TrimSpace(line[5])
		newDetails.Venue = strings.TrimSpace(line[6])
		newDetails.Round = strings.TrimSpace(line[7])
		newDetails.Home = strings.TrimSpace(line[8])
		newDetails.Away = strings.TrimSpace(line[9])
		newResult.Winner = strings.TrimSpace(line[10])
		newResult.Win_by_runs, err = strconv.Atoi(strings.TrimSpace(line[11]))
		newResult.Win_by_wickets, err = strconv.Atoi(strings.TrimSpace(line[12]))
		newResult.Ball_rem, err = strconv.Atoi(strings.TrimSpace(line[13]))
		newDetails.Innings1 = strings.TrimSpace(line[14])
		newScore.Innings1_score, err = strconv.Atoi(strings.TrimSpace(line[15]))
		newScore.Innings1_wickets, err = strconv.Atoi(strings.TrimSpace(line[16]))
		newScore.Innings1_overs_batted, err = strconv.ParseFloat(strings.TrimSpace(line[17]), 64)
		newScore.Innings1_overs, err = strconv.ParseFloat(strings.TrimSpace(line[18]), 64)
		newDetails.Innings2 = strings.TrimSpace(line[19])
		newScore.Innings2_score, err = strconv.Atoi(strings.TrimSpace(line[20]))
		newScore.Innings2_wickets, err = strconv.Atoi(strings.TrimSpace(line[21]))
		newScore.Innings2_overs_batted, err = strconv.ParseFloat(strings.TrimSpace(line[22]), 64)
		newScore.Innings2_overs, err = strconv.ParseFloat(strings.TrimSpace(line[23]), 64)
		newScore.Dl_method, err = strconv.Atoi(strings.TrimSpace(line[24]))
		newScore.Target, err = strconv.Atoi(strings.TrimSpace(line[25]))
		newMatch.Score = newScore
		newMatch.Details = newDetails
		newMatch.Result = newResult
		match = append(match, newMatch)

		index := searchSeries(ser)

		if index != -1 {
			series[index].Matches = append(series[index].Matches, newMatch)
		}
	}
}
func buildSeriesDB() {
	f, err := os.Open("/Users/rishikeshmisal/Work/src/T20-Database-Analytics/t20_series.csv")
	reader := csv.NewReader(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer f.Close()
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		newSeries := Entity.Series{}

		newSeries.Name = line[0]
		newSeries.Season = line[1]
		newSeries.Winner = line[2]
		newSeries.Margin = line[3]
		newSeries.Series_id, err = strconv.Atoi(line[4])

		series = append(series, newSeries)
	}
}
func searchSeries(i int) int {

	for ret, ser := range series {
		if ser.Series_id == i {
			return ret
		}
	}

	return -1
}
func GetMatchesByVenue(writer http.ResponseWriter, request *http.Request) {

	writer.Header().Add("Content-Type", "application/json")
	matchVar := mux.Vars(request)
	venue := matchVar["venue"]
	jsonEncode := json.NewEncoder(writer)
	venue = strings.ToUpper(venue)

	var matches []Entity.Match
	for _, m := range match {
		if strings.Contains(strings.ToUpper(m.Details.Venue), venue) {
			matches = append(matches, m)
		}
	}
	if len(matches) == 0 {
		writer.Write([]byte("Wrong input bhai!"))
		return
	}
	jsonEncode.Encode(matches)
	return

}
func GetTeamMatches(writer http.ResponseWriter, request *http.Request) {
	matchVar := mux.Vars(request)
	team := matchVar["team"]
	jsonEncode := json.NewEncoder(writer)
	team = strings.ToUpper(team)

	writer.Header().Add("Content-Type", "application/json")
	var matches []Entity.Match
	for _, m := range match {
		if strings.ToUpper(m.Details.Innings1) == team || strings.ToUpper(m.Details.Innings2) == team {
			matches = append(matches, m)
		}
	}
	if len(matches) == 0 {
		writer.Write([]byte("Wrong input bhai!"))
		return
	}
	jsonEncode.Encode(matches)
	return

}
func GetSeriesById(writer http.ResponseWriter, request *http.Request) {
	matchVar := mux.Vars(request)
	seriesId, err := strconv.Atoi(matchVar["seriesid"])
	jsonEncode := json.NewEncoder(writer)
	writer.Header().Add("Content-Type", "application/json")
	if err != nil {
		fmt.Printf("Wrong input bhai!, %s", err)

		writer.Write([]byte("Wrong input bhai!"))
		return
	}

	ser := Entity.Series{}
	for _, m := range series {
		if m.Series_id == seriesId {
			ser = m
			break
		}
	}
	jsonEncode.Encode(ser)
	return
}
func GetMatchById(writer http.ResponseWriter, request *http.Request) {

	writer.Header().Add("Content-Type", "application/json")
	matchVar := mux.Vars(request)
	matchId, err := strconv.Atoi(matchVar["matchid"])
	jsonEncode := json.NewEncoder(writer)
	if err != nil {
		fmt.Printf("Wrong input bhai!, %s", err)
	}
	for _, m := range match {
		if m.Id == matchId {
			jsonEncode.Encode(m)
			return
		}
	}
	return

}
