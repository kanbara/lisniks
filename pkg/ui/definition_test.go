package ui

import (
	"testing"
)

func Test_breakWordBoundaries(t *testing.T) {
	type args struct {
		word string
		x    int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test case short word",
			args: args{
				word: "this is a short word",
				x:    50,
			},
			want: "this is a short word",
		},
		{
			name: "test case longer word one split",
			args: args{
				word: "this is a short word",
				x: 14,
			},
			want: "this is a\nshort word",
		},
		{
			name: "test case longer word two split",
			args: args{
				word: "this is a short word",
				x: 9,
			},
			want: "this is\na short\nword",
		},
		{
			name: "test case longer word with newlines already",
			args: args{
				word: "this is a\n\nshort word",
				x: 9,
			},
			want: "this is\na\n\nshort\nword",
		},
		{
			name: "test case longer word with newlines already",
			args: args{
				word: "1. year-old\n\n(1) Used as the second element of compounds with cardinal numbers (in their nominative form, feminine where applicable) to form adjectives describing people's ages.\n\nSi ist tvalif-ašňa talla/tvalif-ašňo. \"She is a twelve-year-old [girl]\". Some numerals have irregular forms:\n\n1: enasans, /e'na:zans/, no secondary stress on the numeral element\n2: tvojasans\n3: řijasans\n4. fidrasans\n5. fimf-asans (the /f/ from the oblique forms reappears)\n\nThere is full assimilation of the palatalisation in spelling, hence: asans ~ ašňa, /a:zans/ ~ /ažňa/",
				x: 62,
			},
			want: "1. year-old\n\n(1) Used as the second element of compounds with cardinal\nnumbers (in their nominative form, feminine where applicable)\nto form adjectives describing people's ages.\n\nSi ist tvalif-ašňa talla/tvalif-ašňo. \"She is a\ntwelve-year-old [girl]\". Some numerals have irregular forms:\n\n1: enasans, /e'na:zans/, no secondary stress on the numeral\nelement\n2: tvojasans\n3: řijasans\n4. fidrasans\n5. fimf-asans (the /f/ from the oblique forms reappears)\n\nThere is full assimilation of the palatalisation in spelling,\nhence: asans ~ ašňa, /a:zans/ ~ /ažňa/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := wordWrap(tt.args.word, tt.args.x); got != tt.want {
				t.Errorf("wordWrap() = `%q`, want `%q`", got, tt.want)
			}
		})
	}
}