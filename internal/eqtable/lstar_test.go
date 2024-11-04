package eqtable

import (
	"github.com/BaldiSlayer/rofl-lab2/internal/defaults"
	"github.com/BaldiSlayer/rofl-lab2/internal/wautomata"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLStar_ToDFA(t *testing.T) {
	type fields struct {
		prefixes []string
		suffixes []string
		answers  []string
		alphabet []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   *wautomata.DFA
	}{
		{
			name: "small",
			fields: fields{
				prefixes: []string{
					"",
					"b",
					"a",
					"ba",
					"bb",
				},
				suffixes: []string{
					"",
				},
				answers: []string{
					"-",
					"+",
					"-",
					"-",
					"-",
				},
				alphabet: []byte{'a', 'b'},
			},
			want: wautomata.New(
				"",
				map[string]struct{}{
					"b": {},
				},
				defaults.GetAlphabet(),
				wautomata.NewTransitions(
					map[wautomata.Transition]string{
						wautomata.Transition{Src: "", Symbol: 'a'}:  "",
						wautomata.Transition{Src: "", Symbol: 'b'}:  "b",
						wautomata.Transition{Src: "b", Symbol: 'a'}: "",
						wautomata.Transition{Src: "b", Symbol: 'b'}: "",
					},
				),
				map[string]struct{}{
					"b": {},
					"":  {},
				},
			),
		},
		{
			name: "from consult",
			fields: fields{
				prefixes: []string{
					"",
					"b",
					"a",
					"ba",
					"bb",
					"aa",
					"ab",
					"baa",
					"bab",
					"bba",
					"bbb",
				},
				suffixes: []string{
					"",
					"a",
					"ba",
					"aba",
				},
				answers: []string{
					"---+",
					"+---",
					"--+-",
					"-+--",
					"----",
					"---+",
					"-+--",
					"+---",
					"----",
					"----",
					"----",
				},
				alphabet: []byte{'a', 'b'},
			},
			want: wautomata.New(
				"",
				map[string]struct{}{
					"b": {},
				},
				defaults.GetAlphabet(),
				wautomata.NewTransitions(
					map[wautomata.Transition]string{
						wautomata.Transition{Src: "", Symbol: 'a'}:   "a",
						wautomata.Transition{Src: "", Symbol: 'b'}:   "b",
						wautomata.Transition{Src: "a", Symbol: 'a'}:  "",
						wautomata.Transition{Src: "a", Symbol: 'b'}:  "ba",
						wautomata.Transition{Src: "b", Symbol: 'a'}:  "ba",
						wautomata.Transition{Src: "b", Symbol: 'b'}:  "bb",
						wautomata.Transition{Src: "ba", Symbol: 'a'}: "b",
						wautomata.Transition{Src: "ba", Symbol: 'b'}: "bb",
						wautomata.Transition{Src: "bb", Symbol: 'a'}: "bb",
						wautomata.Transition{Src: "bb", Symbol: 'b'}: "bb",
					},
				),
				map[string]struct{}{
					"":   {},
					"a":  {},
					"b":  {},
					"ba": {},
					"bb": {},
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alphabet = tt.fields.alphabet

			lstar := &LStar{
				prefixes: tt.fields.prefixes,
				suffixes: tt.fields.suffixes,
				answers:  tt.fields.answers,
			}

			got := lstar.ToDFA()
			require.Equal(t, tt.want, got)
		})
	}
}
