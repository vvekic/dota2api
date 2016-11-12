package dota2api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type pickBan struct {
	IsPick bool `json:"is_pick"`
	HeroID int  `json:"hero_id"`
	Team   int  `json:"team"`
	Order  int  `json:"order"`
}

type abilityUpgrade struct {
	Ability int `json:"ability"`
	Time    int `json:"time"`
	Level   int `json:"level"`
}

type player struct {
	AccountID       int               `json:"account_id"`
	PlayerSlot      int               `json:"player_slot"`
	HeroID          int               `json:"hero_id"`
	Item0           int               `json:"item_0"`
	Item1           int               `json:"item_1"`
	Item2           int               `json:"item_2"`
	Item3           int               `json:"item_3"`
	Item4           int               `json:"item_4"`
	Item5           int               `json:"item_5"`
	Kills           int               `json:"kills"`
	Deaths          int               `json:"deaths"`
	Assists         int               `json:"assists"`
	LeaverStatus    int               `json:"leaver_status"`
	LastHits        int               `json:"last_hits"`
	Denies          int               `json:"denies"`
	GoldPerMin      int               `json:"gold_per_min"`
	XPPerMin        int               `json:"xp_per_min"`
	Level           int               `json:"level"`
	Gold            int               `json:"gold"`
	GoldSpent       int               `json:"gold_spent"`
	HeroDamage      int               `json:"hero_damage"`
	TowerDamage     int               `json:"tower_damage"`
	HeroHealing     int               `json:"hero_healing"`
	AbilityUpgrades []*abilityUpgrade `json:"ability_upgrades"`
}

type match struct {
	Players               []*player     `json:"players"`
	RadiantWin            bool          `json:"radiant_win"`
	Duration              time.Duration `json:"duration"`
	StartTime             int           `json:"start_time"`
	MatchID               int           `json:"match_id"`
	MatchSeqNum           int           `json:"match_seq_num"`
	TowerStatusRadiant    int           `json:"tower_status_radiant"`
	TowerStatusDire       int           `json:"tower_status_dire"`
	BarracksStatusRadiant int           `json:"barracks_status_radiant"`
	BarracksStatusDire    int           `json:"barracks_status_dire"`
	Cluster               int           `json:"cluster"`
	FirstBloodTime        int           `json:"first_blood_time"`
	LobbyType             int           `json:"lobby_type"`
	HumanPlayers          int           `json:"human_players"`
	LeagueID              int           `json:"leagueid"`
	GameMode              int           `json:"game_mode"`
	RadiantScore          int           `json:"radiant_score"`
	DireScore             int           `json:"dire_score"`
	PicksBans             []*pickBan    `json:"picks_bans"`
}

type result struct {
	Status       int      `json:"status"`
	StatusDetail string   `json:"statusDetail"`
	Matches      []*match `json:"matches"`
}

// GetMatchHistoryBySequenceNumResponse contains the matches requested by GetMatchHistoryBySequenceNum
type GetMatchHistoryBySequenceNumResponse struct {
	Result *result `json:"result"`
}

// GetMatchHistoryResponse ...
type GetMatchHistoryResponse struct {
	Result *result `json:"result"`
}

// GetMatchHistoryBySequenceNum returns a slice of Dota 2 matches in the order they were recorded
func (c *Client) GetMatchHistoryBySequenceNum(startAtMatchSeqNum, matchesRequested int) (*GetMatchHistoryBySequenceNumResponse, error) {
	req, err := http.NewRequest(http.MethodGet, "http://api.steampowered.com/IDOTA2Match_570/GetMatchHistoryBySequenceNum/v1/", nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create GetMatchHistoryBySequenceNum request")
	}
	q := req.URL.Query()
	q.Set("key", c.key)
	q.Set("start_at_match_seq_num", fmt.Sprintf("%d", startAtMatchSeqNum))
	req.URL.RawQuery = q.Encode()
	res, err := c.hc.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "GetMatchHistoryBySequenceNum request failed")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetMatchHistoryBySequenceNum HTTP status code %d", res.StatusCode)
	}
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	response := &GetMatchHistoryBySequenceNumResponse{}
	if err := json.Unmarshal(buf, &response); err != nil {
		return nil, err
	}
	if response.Result.Status != 1 {
		return nil, fmt.Errorf("GetMatchHistoryBySequenceNum status %d", response.Result.Status)
	}
	return response, nil
}
