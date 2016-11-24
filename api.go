package dota2api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// PickBan represents a hero pick or ban in the match, if the game mode is Captain's mode
type PickBan struct {
	IsPick bool `json:"is_pick"`
	HeroID int  `json:"hero_id"`
	Team   int  `json:"team"`
	Order  int  `json:"order"`
}

// AbilityUpgrade represents a player's ability upgrades during a match
type AbilityUpgrade struct {
	Ability int `json:"ability"`
	Time    int `json:"time"`
	Level   int `json:"level"`
}

// Player represents a player in a Dota 2 match
type Player struct {
	AccountID         int               `json:"account_id"`
	PlayerSlot        int               `json:"player_slot"`
	HeroID            int               `json:"hero_id"`
	Item0             int               `json:"item_0"`
	Item1             int               `json:"item_1"`
	Item2             int               `json:"item_2"`
	Item3             int               `json:"item_3"`
	Item4             int               `json:"item_4"`
	Item5             int               `json:"item_5"`
	Kills             int               `json:"kills"`
	Deaths            int               `json:"deaths"`
	Assists           int               `json:"assists"`
	LeaverStatus      int               `json:"leaver_status"`
	LastHits          int               `json:"last_hits"`
	Denies            int               `json:"denies"`
	GoldPerMin        int               `json:"gold_per_min"`
	XPPerMin          int               `json:"xp_per_min"`
	Level             int               `json:"level"`
	Gold              int               `json:"gold"`
	GoldSpent         int               `json:"gold_spent"`
	HeroDamage        int               `json:"hero_damage"`
	TowerDamage       int               `json:"tower_damage"`
	HeroHealing       int               `json:"hero_healing"`
	ScaledHeroDamage  int               `json:"scaled_hero_damage"`
	ScaledTowerDamage int               `json:"scaled_tower_damage"`
	ScaledHeroHealing int               `json:"scaled_hero_healing"`
	AbilityUpgrades   []*AbilityUpgrade `json:"ability_upgrades"`
}

// Match represents a Dota 2 match
type Match struct {
	Players               []*Player  `json:"players"`
	RadiantWin            bool       `json:"radiant_win"`
	Duration              int        `json:"duration"`
	StartTime             int        `json:"start_time"`
	MatchID               int        `json:"match_id"`
	MatchSeqNum           int        `json:"match_seq_num"`
	TowerStatusRadiant    int        `json:"tower_status_radiant"`
	TowerStatusDire       int        `json:"tower_status_dire"`
	BarracksStatusRadiant int        `json:"barracks_status_radiant"`
	BarracksStatusDire    int        `json:"barracks_status_dire"`
	Cluster               int        `json:"cluster"`
	FirstBloodTime        int        `json:"first_blood_time"`
	LobbyType             int        `json:"lobby_type"`
	HumanPlayers          int        `json:"human_players"`
	LeagueID              int        `json:"leagueid"`
	GameMode              int        `json:"game_mode"`
	RadiantScore          int        `json:"radiant_score"`
	DireScore             int        `json:"dire_score"`
	PicksBans             []*PickBan `json:"picks_bans"`
}

// Result is the base response structure of the Dota 2 API
type Result struct {
	Status       int      `json:"status"`
	StatusDetail string   `json:"statusDetail"`
	Matches      []*Match `json:"matches"`
}

// GetMatchHistoryBySequenceNumResponse contains the matches requested by GetMatchHistoryBySequenceNum
type GetMatchHistoryBySequenceNumResponse struct {
	Result *Result `json:"result"`
}

// GetMatchHistoryResponse ...
type GetMatchHistoryResponse struct {
	Result *Result `json:"result"`
}

type apiError struct {
	error
}

func (e apiError) TooManyRequests() {}

// GetMatchHistoryBySequenceNum returns a Dota 2 matches in the order they were recorded
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
		if res.StatusCode == http.StatusTooManyRequests {
			return nil, apiError{err}
		}
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

// GetMatchHistory returns Dota 2 matches
func (c *Client) GetMatchHistory() (*GetMatchHistoryResponse, error) {
	req, err := http.NewRequest(http.MethodGet, "http://api.steampowered.com/IDOTA2Match_570/GetMatchHistory/v1/", nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create GetMatchHistory request")
	}
	q := req.URL.Query()
	q.Set("key", c.key)
	req.URL.RawQuery = q.Encode()
	res, err := c.hc.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "GetMatchHistory request failed")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetMatchHistory HTTP status code %d", res.StatusCode)
	}
	if res.StatusCode == http.StatusTooManyRequests {
		return nil, apiError{err}
	}
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	response := &GetMatchHistoryResponse{}
	if err := json.Unmarshal(buf, &response); err != nil {
		return nil, err
	}
	if response.Result.Status != 1 {
		return nil, fmt.Errorf("GetMatchHistory status %d", response.Result.Status)
	}
	return response, nil
}
