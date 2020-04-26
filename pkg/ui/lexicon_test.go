package ui

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"reflect"
	"testing"
)

func Test_updateViewCursorOriginAndState(t *testing.T) {
	type args struct {
		c        coords
		updown   int
		selected int
		minY     int
		maxY     int
	}

	tests := []struct {
		name  string
		args  args
		want  coords
		selected int
	}{
		{
			// ------ ++++++
			// 0 <    a @
			// 1 *	  b
			// 2	  c
			// 3	  d
			// ------ ++++++
			// 4	  e
			// 5	  f $
			name: "test scrolling inside of view from start",
			args: args{
				c: coords{
					// disregard all x coords (for now?)
					cursorPos:   coordinates{y: 0},
					originStart: coordinates{y: 0}, // +++++
					viewSize:    coordinates{y: 4}, // -----
				},
				updown:   1, // *
				selected: 0, // <
				minY:     0, // @
				maxY:     5, // $
			},
			want: coords{
				cursorPos:   coordinates{y: 1},
				originStart: coordinates{y: 0},
			},
			selected: 1,
		},
		{
			// ------ ++++++
			// 0      a @
			// 1 <	  b $
			// 2 *
			// 3
			// ------ ++++++
			name: "test scrolling past small view end",
			args: args{
				c: coords{
					// disregard all x coords (for now?)
					cursorPos:   coordinates{y: 1},
					originStart: coordinates{y: 0}, // +++++
					viewSize:    coordinates{y: 4}, // -----
				},
				updown:   1, // *
				selected: 0, // <
				minY:     0, // @
				maxY:     1, // $
			},
			want: coords{
				cursorPos:   coordinates{y: 1},
				originStart: coordinates{y: 0},
			},
			selected: 1,
		},
		{
			// bug taken from prod
			// the issue was that we had a very small view, but large enough s.t. the
			// updown + cursor.y would be greater than the view.y
			// however the maxY was smaller than view.y so instead of jumping out
			// we should have just did the default case and adjusted the selection
			// accordingly
			name: "test scrolling past small view when there's one frame but we would jump out",
			args: args{
				c: coords{
					cursorPos:   coordinates{y: 43},
					originStart: coordinates{y: 0}, // +++++
					viewSize:    coordinates{y: 52}, // -----
				},
				updown:   22, // *
				selected: 43, // <
				minY:     0, // @
				maxY:     43, // $
			},
			want: coords{
				cursorPos:   coordinates{y: 43},
				originStart: coordinates{y: 0},
			},
			selected: 43,
		},
		{
			// ------ ++++++
			// 0 *    a @
			// 1 <	  b
			// 2	  c
			// 3	  d
			// ------ ++++++
			// 4	  e
			// 5	  f $
			name: "test scrolling to the top inside the view by one",
			args: args{
				c: coords{
					// disregard all x coords (for now?)
					cursorPos:   coordinates{y: 1},
					originStart: coordinates{y: 0}, // +++++
					viewSize:    coordinates{y: 4}, // -----
				},
				updown:   -1, // *
				selected: 1, // <
				minY:     0, // @
				maxY:     5, // $
			},
			want: coords{
				cursorPos:   coordinates{y: 0},
				originStart: coordinates{y: 0},
			},
			selected: 0,
		},
		{
			// ------ ++++++
			// 0 *    a @
			// 1 	  b
			// 2 <	  c
			// 3 	  d
			// ------ ++++++
			// 4	  e
			// 5	  f $
			name: "test scrolling to the top inside the view by several",
			args: args{
				c: coords{
					// disregard all x coords (for now?)
					cursorPos:   coordinates{y: 2},
					originStart: coordinates{y: 0}, // +++++
					viewSize:    coordinates{y: 4}, // -----
				},
				updown:   -3, // *
				selected: 2, // <
				minY:     0, // @
				maxY:     5, // $
			},
			want: coords{
				cursorPos:   coordinates{y: 0},
				originStart: coordinates{y: 0},
			},
			selected: 0,
		},
		{
			// 0	  a @
			// 1	  b
			// ------ ++++++
			// 2      c
			// 3 	  d
			// 4 <	  e
			// 5 *	  f $
			// ------ ++++++
			name: "test scrolling to the bottom inside the view by one",
			args: args{
				c: coords{
					// disregard all x coords (for now?)
					cursorPos:   coordinates{y: 2},
					originStart: coordinates{y: 2}, // +++++
					viewSize:    coordinates{y: 4}, // -----
				},
				updown:   1, // *
				selected: 4, // <
				minY:     0, // @
				maxY:     5, // $
			},
			want: coords{
				cursorPos:   coordinates{y: 3},
				originStart: coordinates{y: 2},
			},
			selected: 5,
		},
		{
			// cursor moves down
			// ------ ++++++
			// 0      a @
			// 1 *	  b
			// 2 <    c
			// 3	  d
			// ------ ++++++
			// 4	  e
			// 5	  f $
			name: "test scrolling inside of view up from middle",
			args: args{
				c: coords{
					// disregard all x coords (for now?)
					cursorPos:   coordinates{y: 2},
					originStart: coordinates{y: 0}, // +++++
					viewSize:    coordinates{y: 4}, // -----
				},
				updown:   -1, // *
				selected: 2, // <
				minY:     0, // @
				maxY:     5, // $
			},
			want: coords{
				cursorPos:   coordinates{y: 1},
				originStart: coordinates{y: 0},
			},
			selected: 1,
		},
		{
			// cursor moves down
			// ------ ++++++
			// 0      a @
			// 1 	  b
			// 2 <	  c
			// 3 *	  d
			// ------ ++++++
			// 4	  e
			// 5	  f $
			name: "test scrolling inside of view down from middle",
			args: args{
				c: coords{
					// disregard all x coords (for now?)
					cursorPos:   coordinates{y: 2},
					originStart: coordinates{y: 0}, // +++++
					viewSize:    coordinates{y: 4}, // -----
				},
				updown:   1, // *
				selected: 2, // <
				minY:     0, // @
				maxY:     5, // $
			},
			want: coords{
				cursorPos:   coordinates{y: 3},
				originStart: coordinates{y: 0},
			},
			selected: 3,
		},
		{
			// origin moves down
			// before        | after a
			// ------ ++++++ | ----- +++++
			// 0      a @	 | 0	 b
			// 1  	  b	  	 | 1	 c
			// 2 	  c	     | 2	 d
			// 3 <	  d	     | 3 *	 e <
			// ------ ++++++ | ----- +++++
			// 4 *	  e	     | 4	 f
			// 5	  f $	 | 5
			name: "test scrolling down one",
			args: args{
				c: coords{
					// disregard all x coords (for now?)
					cursorPos:   coordinates{y: 3},
					originStart: coordinates{y: 0}, // +++++
					viewSize:    coordinates{y: 4}, // -----
				},
				updown:   1, // *
				selected: 3, // <
				minY:     0, // @
				maxY:     5, // $
			},
			want: coords{
				cursorPos:   coordinates{y: 3},
				originStart: coordinates{y: 1},
			},
			selected: 4,
		},
		{
			// origin moves down
			//						 a   @
			// before        | after b
			// ------ ++++++ | ----- +++++
			// 0      a @	 | 0	 c
			// 1  	  b	  	 | 1	 d
			// 2 <	  c	     | 2 <	 e *
			// 3 	  d	     | 3	 f   $
			// ------ ++++++ | ----- +++++
			// 4 *	  e	     | 4
			// 5	  f $	 | 5
			name: "test scrolling down several",
			args: args{
				c: coords{
					// disregard all x coords (for now?)
					cursorPos:   coordinates{y: 2},
					originStart: coordinates{y: 0}, // +++++
					viewSize:    coordinates{y: 4}, // -----
				},
				updown:   2, // *
				selected: 2, // <
				minY:     0, // @
				maxY:     5, // $
			},
			want: coords{
				cursorPos:   coordinates{y: 2},
				originStart: coordinates{y: 2},
			},
			selected: 4,
		},
		{
			// origin moves up
			// before a * @  | after
			// ------ ++++++ | ----- +++++
			// 1 <    b 	 | 0 <	 a * @
			// 2  	  c	 	 | 1	 b
			// 3 	  d	     | 2 	 c
			// 4 	  e	     | 3	 d
			// ------ ++++++ | ----- +++++
			// 5 	  f $    | 4	 e
			//   	   		 | 5	 f   $
			name: "test scrolling up one",
			args: args{
				c: coords{
					// disregard all x coords (for now?)
					cursorPos:   coordinates{y: 0},
					originStart: coordinates{y: 1}, // +++++
					viewSize:    coordinates{y: 4}, // -----
				},
				updown:   -1, // *
				selected: 1, // <
				minY:     0, // @
				maxY:     5, // $
			},
			want: coords{
				cursorPos:   coordinates{y: 0},
				originStart: coordinates{y: 0},
			},
			selected: 0,
		},
		{
			// before  *     | after
			// ------ ++++++ | ----- +++++
			// 0 <    a @	 | 0 <	 a * @
			// 1	  b	  	 | 1	 b
			// 2	  c		 | 2	 c
			// 3	  d		 | 3	 d
			// ------ ++++++ | ----- +++++
			// 4	  e		 | 4	 e
			// 5	  f $	 | 5	 f   $
			name: "test scrolling up past the view by one",
			args: args{
				c: coords{
					// disregard all x coords (for now?)
					cursorPos:   coordinates{y: 0},
					originStart: coordinates{y: 0}, // +++++
					viewSize:    coordinates{y: 4}, // -----
				},
				updown:   -1, // *
				selected: 0, // <
				minY:     0, // @
				maxY:     5, // $
			},
			want: coords{
				cursorPos:   coordinates{y: 0},
				originStart: coordinates{y: 0},
			},
			selected: 0,
		},
		{
			// 0	  a @	   0	 a   @
			// before b      | after b
			// ------ ++++++ | ----- +++++
			// 2      c 	 | 2 	 c
			// 3	  d	  	 | 3	 d
			// 4	  e		 | 4	 e
			// 5 <	  f	$	 | 5 <	 f * $
			// ------ ++++++ | ----- +++++
			//   *
			name: "test scrolling down past the view by one",
			args: args{
				c: coords{
					// disregard all x coords (for now?)
					cursorPos:   coordinates{y: 3},
					originStart: coordinates{y: 2}, // +++++
					viewSize:    coordinates{y: 4}, // -----
				},
				updown:   1, // *
				selected: 5, // <
				minY:     0, // @
				maxY:     5, // $
			},
			want: coords{
				cursorPos:   coordinates{y: 3},
				originStart: coordinates{y: 2},
			},
			selected: 5,
		},
		{
			//						 a   @
			// before        | after b
			// ------ ++++++ | ----- +++++
			// 0      a 	 | 2 	 c
			// 1 	  b	  	 | 3	 d
			// 2	  c		 | 4	 e
			// 3 <	  d	$	 | 5 <	 f * $
			// ------ ++++++ | ----- +++++
			// 4	  e
			// 5 	  f
			//
			//   *
			name: "test scrolling down at end of view doesn't leave empty lines",
			args: args{
				c: coords{
					// disregard all x coords (for now?)
					cursorPos:   coordinates{y: 3},
					originStart: coordinates{y: 0}, // +++++
					viewSize:    coordinates{y: 4}, // -----
				},
				updown:   4, // *
				selected: 3, // <
				minY:     0, // @
				maxY:     5, // $
			},
			want: coords{
				cursorPos:   coordinates{y: 3},
				originStart: coordinates{y: 2},
			},
			selected: 5,
		},
		{
		// before a * @  | after
		// ------ ++++++ | ----- +++++
		// 1 0    b 	 | 0 <	 a * @
		// 2  	  c	 	 | 1	 b
		// 3 	  d	     | 2 	 c
		// 4 < 	  e	     | 3	 d
		// ------ ++++++ | ----- +++++
		// 5 	  f $    | 4	 e
		//   	   		 | 5	 f   $
		name: "test scrolling up at the end doesn't go past 0 either",
		args: args{
			c: coords{
				// disregard all x coords (for now?)
				cursorPos:   coordinates{y: 1},
				originStart: coordinates{y: 1}, // +++++
				viewSize:    coordinates{y: 4}, // -----
			},
			updown:   -7, // *
			selected: 4, // <
			minY:     0, // @
			maxY:     5, // $
		},
		want: coords{
			cursorPos:   coordinates{y: 0},
			originStart: coordinates{y: 0},
		},
		selected: 0,
	},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.name)
			log.SetLevel(log.DebugLevel)
			got, got1 := calculateNewViewAndState(tt.args.c, tt.args.updown, tt.args.selected, tt.args.minY, tt.args.maxY)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calculateNewViewAndState()\n" +
					" got =  cursor(%v) origin(%v)\n" +
					" want = cursor(%v) origin(%v)\n",
					 got.cursorPos.y, got.originStart.y,
					tt.want.cursorPos.y, tt.want.originStart.y)
			}

			if got1 != tt.selected {
				t.Errorf("calculateNewViewAndState() got1 = %v, want %v", got1, tt.selected)
			}
		})
	}
}