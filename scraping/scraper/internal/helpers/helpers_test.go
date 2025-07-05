package helpers

import "testing"

func TestToSnakeCase(t *testing.T) {
	tests := [][2]string{
		{"Main event: Group stage", "main_event:_group_stage"},
		{"Playoff: Semifinal", "playoff:_semifinal"},
		{"Playoff: Grand final", "playoff:_grand_final"},
	}

	for _, test := range tests {
		result := ToSnakeCase(test[0])
		if result != test[1] {
			t.Errorf("Wrong result, want %s, get %s", test[1], result)
		}
	}
}
