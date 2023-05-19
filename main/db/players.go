package db

type Prometheus struct {
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

type Player struct {
	DiscordID      int
	PlayerUsername string
	ELO            int
}
