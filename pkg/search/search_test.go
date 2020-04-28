package search

import (
	"reflect"
	"testing"
)

func TestSearchQueue_Dequeue(t *testing.T) {
	type fields struct {
		data []Data
		max  int
	}

	tests := []struct {
		name      string
		fields    fields
		wantQueue Queue
		wantItem  Data
	}{
		{
			name: "test dequeue item",
			fields: fields{
				data: []Data{
					{Type: 0, Pattern: 0, String: "a"},
					{Type: 0, Pattern: 0, String: "b"},
					{Type: 0, Pattern: 0, String: "c"},
				},
				max: 30,
			},
			wantQueue: Queue{
				data: []Data{
					{Type: 0, Pattern: 0, String: "b"},
					{Type: 0, Pattern: 0, String: "c"},
				},
				max: 30,
			},
			wantItem: Data{Type: 0, Pattern: 0, String: "a"},
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
		data []Data
		max  int
	}

	type args struct {
		str Data
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantData []Data
	}{
		{
			name: "test enqueue with size left",
			fields: fields{
				data: []Data{{Type: 0, Pattern: 0, String: "a"}},
				max:  4,
			},
			args: args{
				str: Data{Type: 0, Pattern: 0, String: "b"},
			},
			wantData: []Data{
				{Type: 0, Pattern: 0, String: "a"},
				{Type: 0, Pattern: 0, String: "b"},
			},
		},
		{
			name: "test enqueue with size left, removing others",
			fields: fields{
				data: []Data{
					{Type: 0, Pattern: 0, String: "a"},
					{Type: 0, Pattern: 0, String: "b"},
					{Type: 0, Pattern: 0, String: "c"},
				},
				max: 5,
			},
			args: args{
				str: Data{Type: 0, Pattern: 0, String: "b"},
			},
			wantData: []Data{
				{Type: 0, Pattern: 0, String: "a"},
				{Type: 0, Pattern: 0, String: "c"},
				{Type: 0, Pattern: 0, String: "b"},
			},
		},
		{
			name: "test enqueue with no size left",
			fields: fields{
				data: []Data{
					{Type: 0, Pattern: 0, String: "a"}, // should end up with {"b", "c"}
					{Type: 0, Pattern: 0, String: "b"},
				},
				max: 2,
			},
			args: args{
				str: Data{Type: 0, Pattern: 0, String: "c"},
			},
			wantData: []Data{
				{Type: 0, Pattern: 0, String: "b"}, // should end up with {"b", "c"}
				{Type: 0, Pattern: 0, String: "c"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewQueue(tt.fields.max)
			s.data = make([]Data, 0, tt.fields.max)
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

//func sptr(str string) *string { return &str }
func TestSearchQueue_Peek(t *testing.T) {
	type fields struct {
		data []Data
		max  int
	}

	type args struct {
		i int
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Data
	}{
		{
			name: "peek valid index",
			fields: fields{
				data: []Data{
					{Type: 0, Pattern: 0, String: "a"},
					{Type: 0, Pattern: 0, String: "b"},
				},
				max: 3,
			},
			args: args{
				i: 1, // goes from right to left
			},
			want: &Data{Type: 0, Pattern: 0, String: "a"},
		},
		{
			name: "peek invalid index",
			fields: fields{
				data: []Data{
					{Type: 0, Pattern: 0, String: "a"},
					{Type: 0, Pattern: 0, String: "b"},
				},
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
				data: []Data{
					{Type: 0, Pattern: 0, String: "a"},
					{Type: 0, Pattern: 0, String: "b"},
				},
				max: 3,
			},
			args: args{
				i: 1,
			},
			want: &Data{Type: 0, Pattern: 0, String: "a"},
		},
		{
			name: "peek first",
			fields: fields{
				data: []Data{
					{Type: 0, Pattern: 0, String: "a"},
					{Type: 0, Pattern: 0, String: "b"},
				},
				max: 3,
			},
			args: args{
				i: 0,
			},
			want: &Data{Type: 0, Pattern: 0, String: "b"},
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
		data []Data
		max  int
	}

	type args struct {
		str Data
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantData []Data
	}{
		{
			name: "test removing others of the same",
			fields: fields{
				data: []Data{
					{Type: 0, Pattern: 0, String: "a"},
					{Type: 0, Pattern: 0, String: "b"},
					{Type: 0, Pattern: 0, String: "c"},
				},
			},
			args: args{
				str: Data{Type: 0, Pattern: 0, String: "b"},
			},
			wantData: []Data{
				{Type: 0, Pattern: 0, String: "a"},
				{Type: 0, Pattern: 0, String: "c"},
			},
		},
		{
			name: "test removing others of the same multiple",
			fields: fields{
				data: []Data{{Type: 0, Pattern: 0, String: "a"},
					{Type: 0, Pattern: 0, String: "b"},
					{Type: 0, Pattern: 0, String: "c"},
					{Type: 0, Pattern: 0, String: "b"},
				},
			},
			args: args{
				str: Data{Type: 0, Pattern: 0, String: "b"},
			},
			wantData: []Data{
				{Type: 0, Pattern: 0, String: "a"},
				{Type: 0, Pattern: 0, String: "c"},
			},
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
