package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
)

const query = `query userContestRankingInfo($userSlug: String!) {
  userContestRanking(userSlug: $userSlug) {
    attendedContestsCount
    rating
    globalRanking
    localRanking
    globalTotalParticipants
    localTotalParticipants
    topPercentage
  }
  userContestRankingHistory(userSlug: $userSlug) {
    attended
    totalProblems
    trendingDirection
    finishTimeInSeconds
    rating
    score
    ranking
    contest {
      title
      titleCn
      startTime
    }
  }
}`

type MyQuery struct {
	Query     string `json:"query"`
	Variables struct {
		UserSlug string `json:"userSlug"`
	} `json:"variables"`
}

type UserRanking struct {
	AttendedContestsCount   int     `json:"attendedContestsCount"`
	Rating                  float64 `json:"rating"`
	GlobalRanking           int     `json:"globalRanking"`
	LocalRanking            int     `json:"localRanking"`
	GlobalTotalParticipants int     `json:"globalTotalParticipants"`
	LocalTotalParticipants  int     `json:"localTotalParticipants"`
	TopPercentage           float64 `json:"topPercentage"`
}

type Contest struct {
	Title     string `json:"title"`
	TitleCn   string `json:"titleCn"`
	StartTime int    `json:"startTime"`
}

type UserHistory struct {
	Attended            bool        `json:"attended"`
	TotalProblems       int         `json:"totalProblems"`
	TrendingDirection   interface{} `json:"trendingDirection"`
	FinishTimeInSeconds int         `json:"finishTimeInSeconds"`
	Rating              float64     `json:"rating"`
	Score               int         `json:"score"`
	Ranking             int         `json:"ranking"`
	Contest             Contest     `json:"contest"`
}

type Response struct {
	Data struct {
		UserContestRanking        UserRanking   `json:"userContestRanking"`
		UserContestRankingHistory []UserHistory `json:"userContestRankingHistory"`
	} `json:"data"`
}

type OneUser struct {
	UserName                  string        `json:"username"`
	UserContestRanking        UserRanking   `json:"userContestRanking"`
	UserContestRankingHistory []UserHistory `json:"userContestRankingHistory"`
}

type Output struct {
	Users []OneUser `json:"users"`
	Code  int       `json:"code"`
}

const LC_URL = "https://leetcode.com/graphql/noj-go/"

func main() {
	var fname, ofile string
	flag.StringVar(&fname, "in", "users.txt", "User File")
	flag.StringVar(&ofile, "out", "users.json", "Output File")
	flag.Parse()

	f, err := os.Open(fname)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	bt, err := ioutil.ReadAll(f)
	if err != nil {
		log.Panic(err)
	}

	users := strings.Split(strings.Trim(string(bt), "\n \t\r"), "\n")

	var output Output
	output.Code = 0

	for _, user := range users {
		fmt.Printf("<%s>\n", user)
		r := AskUser(user)
		output.Users = append(output.Users, r)
	}
	f2, err := os.Create(ofile)
	if err != nil {
		log.Panic(err)
	}
	defer f2.Close()
	bt, err = json.Marshal(output)
	if err != nil {
		log.Panic(err)
	}
	f2.Write(bt)
}

func AskUser(name string) OneUser {
	var q MyQuery
	q.Query = query
	q.Variables.UserSlug = name

	var result Response
	resty.New().R().SetBody(q).SetResult(&result).Post(LC_URL)
	var oneUser OneUser
	oneUser.UserContestRanking = result.Data.UserContestRanking
	oneUser.UserContestRankingHistory = result.Data.UserContestRankingHistory
	oneUser.UserName = name
	return oneUser
}
