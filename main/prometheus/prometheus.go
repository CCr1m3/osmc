package prometheus

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ccr1m3/osmc/main/env"
	"github.com/ccr1m3/osmc/main/static"
)

type PrometheusQueryResponse struct {
	Matches []struct {
		Username    string `json:"username"`
		PlayerID    string `json:"playerId"`
		LogoID      string `json:"logoId"`
		Title       string `json:"title"`
		NameplateID string `json:"nameplateId"`
		EmoticonID  string `json:"emoticonId"`
		TitleID     string `json:"titleId"`
		Tags        []any  `json:"tags"`
		PlatformIds struct {
		} `json:"platformIds"`
		MasteryLevel int `json:"masteryLevel"`
		Organization struct {
			OrganizationID string `json:"organizationId"`
			LogoID         string `json:"logoId"`
			Name           string `json:"name"`
		} `json:"organization"`
	} `json:"matches"`
	Paging struct {
		Page     int `json:"page"`
		PageSize int `json:"pageSize"`
	} `json:"paging"`
}

type PrometheusSearchResponse struct {
	Players []struct {
		Username    string   `json:"username"`
		PlayerID    string   `json:"playerId"`
		LogoID      string   `json:"logoId"`
		Title       string   `json:"title"`
		NameplateID string   `json:"nameplateId"`
		EmoticonID  string   `json:"emoticonId"`
		TitleID     string   `json:"titleId"`
		Tags        []string `json:"tags"`
		PlatformIds struct {
		} `json:"platformIds"`
		MasteryLevel int `json:"masteryLevel"`
		Organization struct {
			OrganizationID string `json:"organizationId"`
			LogoID         string `json:"logoId"`
			Name           string `json:"name"`
		} `json:"organization,omitempty"`
		SocialURL            string `json:"socialUrl,omitempty"`
		Rank                 int    `json:"rank"`
		Wins                 int    `json:"wins"`
		Losses               int    `json:"losses"`
		Games                int    `json:"games"`
		TopRole              string `json:"topRole"`
		Rating               int    `json:"rating"`
		MostPlayedCharacters []struct {
			CharacterID string `json:"characterId"`
			GamesPlayed int    `json:"gamesPlayed"`
		} `json:"mostPlayedCharacters"`
		CurrentDivisionID string `json:"currentDivisionId"`
		ProgressToNext    int    `json:"progressToNext"`
	} `json:"players"`
	Paging struct {
		StartRank  int `json:"startRank"`
		PageSize   int `json:"pageSize"`
		TotalItems int `json:"totalItems"`
	} `json:"paging"`
}

type PlayerResponse struct {
	Username string
	Elo      int
}

func GetRankInfoFromUsername(ctx context.Context, username string) (*PlayerResponse, error) {
	var bearer = "Bearer " + env.Prometheus.Authorization
	var refreshtoken = env.Prometheus.Refreshtoken
	prometheusUrl := fmt.Sprintf("https://prometheus-proxy.odysseyinteractive.gg/api/v1/players?usernameQuery=%s", url.PathEscape(username))
	req, err := http.NewRequest("GET", prometheusUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Authorization", bearer)
	req.Header.Add("X-Refresh-Token", refreshtoken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	jsonBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var queryResponse PrometheusQueryResponse
	err = json.Unmarshal(jsonBytes, &queryResponse)
	if err != nil {
		return nil, err
	}
	var returnInfo PlayerResponse
	if len(queryResponse.Matches) == 1 {
		if strings.EqualFold(username, queryResponse.Matches[0].Username) {
			usernameID := queryResponse.Matches[0].PlayerID
			prometheusUrl = fmt.Sprintf("https://prometheus-proxy.odysseyinteractive.gg/api/v1/ranked/leaderboard/search/%s", url.PathEscape(usernameID))
			req, err := http.NewRequest("GET", prometheusUrl, nil)
			if err != nil {
				return nil, err
			}
			req.Header.Add("X-Authorization", bearer)
			req.Header.Add("X-Refresh-Token", refreshtoken)
			resp, err := client.Do(req)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()
			jsonBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			var searchResponse PrometheusSearchResponse
			err = json.Unmarshal(jsonBytes, &searchResponse)
			if err != nil {
				return nil, err
			}
			for _, player := range searchResponse.Players {
				if player.PlayerID == usernameID {
					returnInfo = PlayerResponse{Username: username, Elo: player.Rating}
					break
				}
			}
			if strings.Compare(returnInfo.Username, "") == 0 {
				returnInfo = PlayerResponse{Username: username, Elo: 0}
				return &returnInfo, static.ErrUsernameNotOnGlobal
			} else {
				return &returnInfo, nil
			}
		} else {
			return nil, static.ErrUsernameNotFound
		}
	} else {
		return nil, static.ErrUsernameNotFound
	}
}
