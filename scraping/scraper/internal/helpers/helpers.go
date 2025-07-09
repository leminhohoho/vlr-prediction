package helpers

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"
)

func ToSnakeCase(str string) string {
	lowerCasedStr := strings.ToLower(str)
	fragments := strings.Fields(lowerCasedStr)
	return strings.Join(fragments, "_")
}

func TimeToSeconds(timeStr string) (int, error) {
	colonCount := strings.Count(timeStr, ":")

	var err error
	var t time.Time

	switch colonCount {
	case 1:
		t, err = time.Parse("15:04:05", "0:"+timeStr)
		if err != nil {
			return -1, err
		}
	case 2:
		t, err = time.Parse("15:04:05", timeStr)
		if err != nil {
			return -1, err
		}
	default:
		return -1, fmt.Errorf("%s is not parsable to duration", timeStr)
	}

	return int(t.Second() + t.Minute()*60 + t.Hour()*3600), nil
}

func FillPlayerKDA(defStat, atkStat, bothSideStat *int) (*int, *int, error) {
	if defStat != nil && atkStat != nil {
		return defStat, atkStat, nil
	}

	if bothSideStat == nil {
		return nil, nil, fmt.Errorf("Can't fill in the missing stat without both side stat")
	}

	if defStat == nil {
		filledDefStat := (*bothSideStat - *atkStat)
		return &filledDefStat, atkStat, nil
	} else {
		filledAtkStat := (*bothSideStat - *defStat)
		return defStat, &filledAtkStat, nil
	}
}

func FillPlayerPerRoundStat(
	defStat, atkStat, bothSideStat *float64,
	teamDefRounds, teamAtkRounds int,
) (*float64, *float64, error) {
	if defStat != nil && atkStat != nil {
		return defStat, atkStat, nil
	}

	if bothSideStat == nil {
		return nil, nil, fmt.Errorf("Can't fill in the missing stat without both side stat")
	}

	totalRounds := float64(teamDefRounds + teamAtkRounds)

	var filledStat float64

	if defStat == nil {
		if teamDefRounds == 0 {
			filledStat = 0
		} else {
			filledStat = (float64(*bothSideStat)*totalRounds - float64(teamAtkRounds)*float64(*atkStat)) / float64(
				teamDefRounds,
			)
		}

		return &filledStat, atkStat, nil
	} else {
		if teamAtkRounds == 0 {
			filledStat = 0
		} else {
			filledStat = (float64(*bothSideStat)*totalRounds - float64(teamDefRounds)*float64(*defStat)) / float64(
				teamAtkRounds,
			)
		}

		return defStat, &filledStat, nil
	}
}

func FillPlayerPerKillStat(
	defStat, atkStat, bothSideStat *float64,
	playerDefKills, playerAtkKills int,
) (*float64, *float64, error) {
	if defStat != nil && atkStat != nil {
		return defStat, atkStat, nil
	}

	if bothSideStat == nil {
		return nil, nil, fmt.Errorf("Can't fill in the missing stat without both side stat")
	}

	totalKills := float64(playerDefKills + playerAtkKills)

	var filledStat float64

	if defStat == nil {
		if playerDefKills == 0 {
			filledStat = 0
		} else {
			filledStat = (float64(*bothSideStat)*totalKills - float64(playerAtkKills)*float64(*atkStat)) / float64(
				playerDefKills,
			)
		}

		return &filledStat, atkStat, nil
	} else {
		if playerAtkKills == 0 {
			filledStat = 0
		} else {
			filledStat = (float64(*bothSideStat)*totalKills - float64(playerDefKills)*float64(*defStat)) / float64(
				playerAtkKills,
			)
		}

		return defStat, &filledStat, nil
	}
}

func GetConn(dbPath string) (*sql.DB, error) {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Database file does not exist: %s", dbPath)
	}

	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
