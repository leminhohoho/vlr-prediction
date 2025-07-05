package urlinfo

import (
	"testing"
)

func TestUrlInfo(t *testing.T) {
	matchUrl := "/506931/bilibili-gaming-vs-tyloo-vct-2025-china-stage-2-w1"
	teamUrl := "/team/12010/bilibili-gaming"
	tournamentUrl := "/event/2499/vct-2025-china-stage-2/group-stage"
	playerUrl := "/player/10698/nephh"

	matchUrlInfo, err := ExtractUrlInfo(matchUrl)
	if err != nil {
		t.Error(err)
	}

	teamUrlInfo, err := ExtractUrlInfo(teamUrl)
	if err != nil {
		t.Error(err)
	}

	tournamentUrlInfo, err := ExtractUrlInfo(tournamentUrl)
	if err != nil {
		t.Error(err)
	}

	playerUrlInfo, err := ExtractUrlInfo(playerUrl)
	if err != nil {
		t.Error(err)
	}

	if !matchUrlInfo.IsMatch() {
		t.Errorf("Match url info has wrong url type")
	}

	if !teamUrlInfo.IsTeam() {
		t.Errorf("Team url info has wrong url type")
	}

	if !playerUrlInfo.IsPlayer() {
		t.Errorf("Player url info has wrong url type")
	}

	if !tournamentUrlInfo.IsTournament() {
		t.Errorf("Tournament url info has wrong url type")
	}

	if matchUrlInfo.Id != 506931 {
		t.Errorf("Wrong match id, want 506931, got %d", matchUrlInfo.Id)
	}

	if teamUrlInfo.Id != 12010 {
		t.Errorf("Wrong team id, want 12010, got %d", teamUrlInfo.Id)
	}

	if tournamentUrlInfo.Id != 2499 {
		t.Errorf("Wrong tournament id, want 2499, got %d", tournamentUrlInfo.Id)
	}

	if playerUrlInfo.Id != 10698 {
		t.Errorf("Wrong player id, want 10698, got %d", playerUrlInfo.Id)
	}
}
