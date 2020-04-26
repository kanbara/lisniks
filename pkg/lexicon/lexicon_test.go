package lexicon

import (
	"github.com/kanbara/lisniks/pkg/language"
	"reflect"
	"testing"
)

func TestService_findConWords(t *testing.T) {
	type fields struct {
		lexicon    Lexicon
		alphaOrder language.AlphaOrderMap
	}

	type args struct {
		str   string
		fuzzy bool
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   Lexicon
	}{
		{
			name: "test find nonfuzzy conword",
			fields: fields{
				lexicon:    Lexicon{
					{Con:     "adgemoǆan"},
					{Con:     "adhærents"},
					{Con:     "ad i"},
					{Con:     "ferbužan"},
					{Con:     "ad"},
				},
			},
			args: args{
				str:   "ad",
				fuzzy: false,
			},
			want: Lexicon{
				{Con:     "adgemoǆan"},
				{Con:     "adhærents"},
				{Con:     "ad i"},
				{Con:     "ad"},
			},
		},
		{
			name: "test find fuzzy conword",
			fields: fields{
				lexicon:    Lexicon{
					{Con:     "adgemoǆan"},
					{Con:     "adhærents"},
					{Con:     "ad i"},
					{Con:     "ferbužan"},
					{Con:     "ad"},
				},
			},
			args: args{
				str:   "a",
				fuzzy: true,
			},
			want: Lexicon{
				{Con:     "adgemoǆan"},
				{Con:     "adhærents"},
				{Con:     "ad i"},
				{Con:     "ferbužan"},
				{Con:     "ad"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				lexicon:    tt.fields.lexicon,
				alphaOrder: tt.fields.alphaOrder,
			}

			if got := s.findConWords(tt.args.str, tt.args.fuzzy); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findConWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func makeAlphaOrderMap(s string) language.AlphaOrderMap {
	runes := []rune(s)
	ao := make(language.AlphaOrderMap, len(runes))

	for i := 0; i < len(runes); i++ {
		ao[runes[i]] = int64(i)
	}

	return ao
}

func TestService_Less(t *testing.T) {
	type fields struct {
		lexicon    Lexicon
		alphaOrder language.AlphaOrderMap
	}

	type args struct {
		i int
		j int
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "test normal alphaorder sorting a < b",
			fields: fields{
				lexicon:    Lexicon{
					{Con: "bužan"},
					{Con: "ad"},
				},
				alphaOrder: makeAlphaOrderMap("abdnuž"),
			},
			args: args{
				i: 0,
				j: 1,
			},
			want: false,
		},
		{
			name: "test normal alphaorder sorting smaller word later",
			fields: fields{
				lexicon:    Lexicon{
					{Con: "ad"},
					{Con: "ad i"},
				},
				alphaOrder: makeAlphaOrderMap("adi"),
			},
			args: args{
				i: 0,
				j: 1,
			},
			want: true,
		},
		{
			name: "test normal alphaorder sorting same length with space goes later",
			fields: fields{
				lexicon:    Lexicon{
					{Con: "ada"},
					{Con: "ad i"},
				},
				alphaOrder: makeAlphaOrderMap("adi"),
			},
			args: args{
				i: 0,
				j: 1,
			},
			want: true,
		},
		{
			name: "test normal alphaorder sorting larger word goes first if the substring is less",
			fields: fields{
				lexicon:    Lexicon{
					{Con: "dusvimman"},
					{Con: "dušks"},
				},
				alphaOrder: makeAlphaOrderMap("abimnsšuv"),
			},
			args: args{
				i: 0,
				j: 1,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				lexicon:    tt.fields.lexicon,
				alphaOrder: tt.fields.alphaOrder,
			}

			if got := s.Less(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("Less() = %v, want %v", got, tt.want)
			}
		})
	}
}