package prometheus

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PrometheusResponse struct {
	Players []struct {
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
		} `json:"organization,omitempty"`
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
		SocialURL         string `json:"socialUrl,omitempty"`
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
	prometheusUrl := fmt.Sprintf("https://prometheus.odysseyinteractive.gg/api/v1/ranked/leaderboard/players?startRank=1&pageSize=100")
	resp, err := http.Get(prometheusUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	jsonBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response PrometheusResponse
	err = json.Unmarshal(jsonBytes, &response)
	if err != nil {
		return nil, err
	}
	var returnInfo PlayerResponse
	returnInfo = PlayerResponse{Username: "", Elo: 0}
	return &returnInfo, nil
}

// if strings.EqualFold(response.Error, "") {
// 	if !response.RankedStats.IsRanked {
// 		corestrikeUrl := fmt.Sprintf("https://corestrike.gg/lookup/%s?region=Europe&json=true", url.PathEscape(username))
// 		resp, err := http.Get(corestrikeUrl)
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer resp.Body.Close()
// 		jsonBytes, err := io.ReadAll(resp.Body)
// 		if err != nil {
// 			return nil, err
// 		}
// 		err = json.Unmarshal(jsonBytes, &response)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	if response.RankedStats.IsRanked {
// 		return &response, nil
// 	} else {
// 		return &response, static.ErrUnrankedUser
// 	}
// } else if strings.EqualFold(response.Error, "Invalid username") {
// 	return nil, static.ErrUsernameInvalid
// } else {
// 	return nil, static.ErrCorestrikeNotFound
// }
