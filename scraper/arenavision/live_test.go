package arenavision

import "testing"

func TestStringToLives(t *testing.T) {
	input := "11-12 [SPA] <br />       13-14-S3 [SRB] 14 [ITA]"
	output := []*Live{
		&Live{"11", "SPA"},
		&Live{"12", "SPA"},
		&Live{"13", "SRB"},
		&Live{"14", "SRB"},
		&Live{"S3", "SRB"},
		&Live{"14", "ITA"},
	}
	res := stringToLives(input)
	if len(res) != len(output) {
		t.Errorf("expecting len=%d, got %d", len(output), len(res))
		return
	}
	for j, r := range res {
		o := output[j]
		if (r.ChannelID != o.ChannelID) || (r.Lang != o.Lang) {
			t.Errorf("item[%d]: expecting %v, got %v", j, o, r)
		}
	}
}
