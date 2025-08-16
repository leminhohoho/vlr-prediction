package helpers

import (
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	tests := [][2]string{
		{"Main event: Group stage", "main_event:_group_stage"},
		{"Playoff: Semifinal", "playoff:_semifinal"},
		{"Playoff: Grand final", "playoff:_grand_final"},
		{"Playoffs: Upper Round 1", "playoffs:_upper_round_1"},
	}

	for _, test := range tests {
		result := ToSnakeCase(test[0])
		if result != test[1] {
			t.Errorf("Wrong result, want %s, get %s", test[1], result)
		}
	}
}

type TimeToSecondsTest struct {
	timeStr  string
	duration int
}

func TestTimeToSeconds(t *testing.T) {
	tests := []TimeToSecondsTest{
		{"00:00", 0},
		{"49:10", 2950},
		{"1:50:20", 6620},
	}

	for _, test := range tests {
		duration, err := TimeToSeconds(test.timeStr)
		if err != nil {
			t.Fatal(err)
		}

		if duration != test.duration {
			t.Errorf("Wrong duration, want %d, get %d instead", test.duration, duration)
		}
	}
}
