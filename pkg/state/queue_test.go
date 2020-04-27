package state

import (
	"reflect"
	"testing"
)

func TestStack_Dequeue(t *testing.T) {
	type fields struct {
		data data
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
				data: []string{"a","b", "c"},
				max:  30,
			},
			wantQueue: Queue{
				data: []string{"b", "c"},
				max:  30,
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

func TestStack_Enqueue(t *testing.T) {
	type fields struct {
		data data
		max  int
	}

	type args struct {
		str string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		wantData data
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
			wantData: []string{"a", "b"},
		},
		{
			name: "test enqueue with size left, removing others",
			fields: fields{
				data: []string{"a", "b", "c"},
				max:  5,
			},
			args: args{
				str: "b",
			},
			wantData: []string{"a", "c", "b"},

		},
		{
			name: "test enqueue with no size left",
			fields: fields{
				data: []string{"a", "b"}, // should end up with {"b", "c"}
				max:  2,
			},
			args: args{
				str: "c",
			},
			wantData: []string{"b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewQueue(tt.fields.max)
			s.data = make(data, 0, tt.fields.max)
			for i := range tt.fields.data {
				s.data = append(s.data, tt.fields.data[i])
			}

			s.Enqueue(tt.args.str)

			if s.data[len(s.data)-1] != tt.args.str {
				t.Errorf("Enqueue() got = %v, want %v", s.data[len(s.data)-1], tt.args.str)
			}

			if !reflect.DeepEqual(s.data, tt.wantData) {
				t.Errorf("data doesn't match: got%#v, want %#v", s.data, tt.wantData)

			}
		})

	}
}

func sptr(str string) *string { return &str }
func TestQueue_Peek(t *testing.T) {
	type fields struct {
		data data
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
				max:  3,
			},
			args: args{
				i: 1, // goes from right to left
			},
			want: sptr("a"),
		},
		{
			name: "peek invalid index",
			fields: fields{
				data: []string{"a","b"},
				max:  3,
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
				max:  3,
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
				max:  3,
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
				t.Errorf("Peek() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestQueue_RemoveOthers(t *testing.T) {
	type fields struct {
		data data
		max  int
	}

	type args struct {
		str string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		wantData data
	}{
		{
			name: "test removing others of the same",
			fields: fields{
				data: []string{"a", "b", "c"},
			},
			args: args{
				str: "b",
			},
			wantData: []string{"a","c"},
		},
		{
			name: "test removing others of the same multiple",
			fields: fields{
				data: []string{"a", "b", "c", "b"},
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