package arenavision

import "testing"

func TestStringToLives(t *testing.T) {
	input := "11-12 [SPA] <br />       13-14-S3 [SRB] 14 [ITA]"
	output := []*Live{
		&Live{Channel("11"), "SPA"},
		&Live{Channel("12"), "SPA"},
		&Live{Channel("13"), "SRB"},
		&Live{Channel("14"), "SRB"},
		&Live{Channel("S3"), "SRB"},
		&Live{Channel("14"), "ITA"},
	}
	res := stringToLives(input)
	if len(res) != len(output) {
		t.Errorf("expecting len=%d, got %d", len(output), len(res))
		return
	}
	for j, r := range res {
		o := output[j]
		if (r.channel != o.channel) || (r.lang != o.lang) {
			t.Errorf("item[%d]: expecting %v, got %v", j, o, r)
		}
	}
}
