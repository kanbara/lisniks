package search

import (
	"reflect"
	"testing"
)

func TestSearchQueue_Dequeue(t *testing.T) {
	type fields struct {
		data []string
		max  int
	}

	tests := []struct {
		name      string
		fields    fields
		wantQueue Queue
		wantItem  string
	}{
		{
			name: "test dequeue item",
			fields: fields{
				data: []string{"a","b","c"},
				max: 30,
			},
			wantQueue: Queue{
				data: []string{"b","c"},
				max: 30,
			},
			wantItem: "a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Queue{
				data: tt.fields.data,
				max:  tt.fields.max,
			}

			item := s.Dequeue()
			if !reflect.DeepEqual(*s, tt.wantQueue) {
				t.Errorf("Dequeue() got stack = %#v, want %#v", *s, tt.wantQueue)
			}

			if *item != tt.wantItem {
				t.Errorf("Dequeue() got item = %v, want %v", *item, tt.wantItem)
			}
		})
	}
}

func TestSearchQueue_Enqueue(t *testing.T) {
	type fields struct {
		data []string
		max  int
	}

	type args struct {
		str string
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantData []string
	}{
		{
			name: "test enqueue with size left",
			fields: fields{
				data: []string{"a"},
				max:  4,
			},
			args: args{
				str: "b",
			},
			wantData: []string{"a","b"},
		},
		{
			name: "test enqueue with size left, removing others",
			fields: fields{
				data: []string{"a","b","c"},
				max: 5,
			},
			args: args{
				str: "b",
			},
			wantData: []string{"a","c","b"},
		},
		{
			name: "test enqueue with no size left",
			fields: fields{
				data: []string{"a","b"},
				max: 2,
			},
			args: args{
				str: "c",
			},
			wantData: []string{"b","c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewQueue(tt.fields.max)
			s.data = make([]string, 0, tt.fields.max)
			for i := range tt.fields.data {
				s.data = append(s.data, tt.fields.data[i])
			}

			s.Enqueue(tt.args.str)

			if s.data[len(s.data)-1] != tt.args.str {
				t.Errorf("Enqueue() got = %v, want %v", s.data[len(s.data)-1], tt.args.str)
			}

			if !reflect.DeepEqual(s.data, tt.wantData) {
				t.Errorf("data doesn't match:\ngot %v\nwant %v", s.data, tt.wantData)

			}
		})

	}
}

func sptr(str string) *string { return &str }

func TestSearchQueue_Peek(t *testing.T) {
	type fields struct {
		data []string
		max  int
	}

	type args struct {
		i int
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   *string
	}{
		{
			name: "peek valid index",
			fields: fields{
				data: []string{"a","b"},
				max: 3,
			},
			args: args{
				i: 1, // goes from right to left
			},
			want: sptr("a"),
		},
		{
			name: "peek invalid index",
			fields: fields{
				data: []string{"a", "b"},
				max: 3,
			},
			args: args{
				i: 7,
			},
			want: nil,
		},
		{
			name: "peek last",
			fields: fields{
				data: []string{"a","b"},
				max: 3,
			},
			args: args{
				i: 1,
			},
			want: sptr("a"),
		},
		{
			name: "peek first",
			fields: fields{
				data: []string{"a","b"},
				max: 3,
			},
			args: args{
				i: 0,
			},
			want: sptr("b"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Queue{
				data: tt.fields.data,
				max:  tt.fields.max,
			}

			if got := s.Peek(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Peek() = %v, want %v", *got, tt.want)
			}
		})
	}
}

func TestSearchQueue_RemoveOthers(t *testing.T) {
	type fields struct {
		data []string
		max  int
	}

	type args struct {
		str string
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantData []string
	}{
		{
			name: "test removing others of the same",
			fields: fields{
				data: []string{"a","b","c"},
			},
			args: args{
				str: "b",
			},
			wantData: []string{"a","c"},
		},
		{
			name: "test removing others of the same multiple",
			fields: fields{
				data: []string{"a","b","c","b"},
			},
			args: args{
				str: "b",
			},
			wantData: []string{"a","c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Queue{
				data: tt.fields.data,
				max:  tt.fields.max,
			}

			s.RemoveOthers(tt.args.str)

			if !reflect.DeepEqual(s.data, tt.wantData) {
				t.Errorf("RemoveOthers() = %v, want %v", s.data, tt.wantData)
			}
		})
	}
}

func Test_parseTypeAndPattern(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name  string
		args  args
		wantErr bool
		t Type
		p Pattern
	}{
		{
			name: "test only type",
			args: args{
				str: "a",
			},
			t:     TypeAustrianWord,
			p:     PatternRegex,
			wantErr: false,
		},
		{
			name: "test type and pattern",
			args: args{
				str: "ap",
			},
			t:       TypeAustrianWord,
			p:       PatternPhonotactic,
			wantErr: false,
		},
		{
			name: "test type and pattern II",
			args: args{
				str: "dp",
			},
			t:       TypeWordDefinition,
			p:       PatternPhonotactic,
			wantErr: false,
		},
		{
			name: "test only pattern",
			args: args{
				str: "p",
			},
			t:       TypeAustrianWord,
			p:       PatternPhonotactic,
			wantErr: false,
		},
		{
			name: "test defaults",
			args: args{
				str: "",
			},
			t:       TypeAustrianWord,
			p:       PatternRegex,
			wantErr: false,
		},
		{
			name: "test setting type twice",
			args: args{
				str: "ae",
			},
			wantErr: true,
		},
		{
			name: "test setting pattern twice",
			args: args{
				str: "rp",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st, pp, err := parseTypeAndPattern(tt.args.str)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseTypeAndPattern() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(st, tt.t) {
				t.Errorf("parseTypeAndPattern() type = %v, want %v", st, tt.t)
			}
			if !reflect.DeepEqual(pp, tt.p) {
				t.Errorf("parseTypeAndPattern() pattern = %v, want %v", pp, tt.p)
			}
		})
	}
}

func TestParseString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    Parsed
		wantErr bool
	}{
		{
			name: "test normal case",
			args: args{
				str: "e/fish",
			},
			want: Parsed{
				Type:    TypeEnglishWord,
				Pattern: PatternRegex,
				String:  "fish",
			},
			wantErr: false,
		},
		{
			name: "test weird stuff too many slashes after",
			args: args{
				str: "e/fish/foobar",
			},
			wantErr: true,
		},
		{
			name: "test weird stuff too many slashes before",
			args: args{
				str: "/e/fish",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseString(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseString() got = %v, want %v", got, tt.want)
			}
		})
	}
}