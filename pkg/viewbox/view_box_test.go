package viewbox

import (
	"testing"
)

func TestSplitStringTD(t *testing.T) {
	vb := ViewBox{}

	group := []struct {
		name,
		input string
		want int
	}{
		{"eagle should be 5", "eagle", 5},
		{"ołówek should be 6", "ołówek", 6},
		{"deuxième should be 8", "deuxième", 8},
		{"Страсбург анттары should be 17", "Страсбург анттары", 17},
		{"მუდანობაშ should be 9", "მუდანობაშ", 9},
	}

	for _, tt := range group {
		t.Run(tt.name, func(t *testing.T) {
			answ, _ := vb.splitString(tt.input)

			if answ != tt.want {
				t.Errorf("Got %d, want %d", answ, tt.want)
			}
		})
	}
}

func TestAlignLeftTD(t *testing.T) {
	vb := ViewBox{}

	group := []struct {
		name,
		st,
		padding,
		want string
		maxWidth int
	}{
		{"1", "Страсбург", " ", "Страсбург", 3},
		{"2", "Страсбург", "~", "Страсбург~~~", 12},
		{"3", "Страсбург", " ", "Страсбург         ", 18},
		{"4", "Страсбург", ".", "Страсбург....", 13},
	}

	for _, tt := range group {
		t.Run(tt.name, func(t *testing.T) {
			answ := vb.AlignLeft(tt.maxWidth, tt.st, tt.padding)

			if answ != tt.want {
				t.Errorf("Got %s, want %s", answ, tt.want)
			}
		})
	}
}

func TestAlignRightTD(t *testing.T) {
	vb := ViewBox{}

	group := []struct {
		name,
		st,
		padding,
		want string
		maxWidth int
	}{
		{"1", "Страсбург", " ", "Страсбург", 3},
		{"2", "Страсбург", "~", "~~~Страсбург", 12},
		{"3", "Страсбург", " ", "         Страсбург", 18},
		{"4", "Страсбург", ".", "....Страсбург", 13},
	}

	for _, tt := range group {
		t.Run(tt.name, func(t *testing.T) {
			answ := vb.AlignRight(tt.maxWidth, tt.st, tt.padding)

			if answ != tt.want {
				t.Errorf("Got %s, want %s", answ, tt.want)
			}
		})
	}
}

func TestBreakString(t *testing.T) {
	vb := ViewBox{}

	testString := "Some not so very long sentence, but long enough to be tested."

	maxWidth := 21

	want := []string{"Some not so very long", " sentence, but long e", "nough to be tested."}

	answ := vb.breakString(testString, maxWidth)

	equal := true

	for k, v := range want {
		if v != answ[k] {
			equal = false
			break
		}
	}

	if !equal {
		t.Errorf("TestWrapString Got '%v', want '%v'", answ[1], want[1])
	}
}

func TestTrimEnd(t *testing.T) {
	vb := ViewBox{}

	input := "Some lengthy string input"

	want := "Some lengthy stri.."

	answ := vb.TrimEnd(input, 19, 19, 2, '.')

	if want != answ {
		t.Errorf("TestTrimEnd, got: %s, want %s", answ, want)
	}
}
