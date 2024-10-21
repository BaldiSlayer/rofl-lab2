package mat

import (
	"testing"

	"github.com/BaldiSlayer/rofl-lab2/internal/maze"
	"github.com/stretchr/testify/require"
)

func TestRealization_Include(t *testing.T) {
	type fields struct {
		maze *maze.Maze
	}

	type args struct {
		query string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "1",
			fields: fields{
				maze: maze.New([][]bool{
					{true, true},
					{true, true},
				}),
			},
			args: args{
				query: "SSS",
			},
			want: true,
		},
		{
			name: "2",
			fields: fields{
				maze: maze.New([][]bool{
					{true, true},
					{false, true},
				}),
			},
			args: args{
				query: "SSS",
			},
			want: false,
		},
		{
			name: "3",
			fields: fields{
				maze: maze.New([][]bool{
					{true, true},
					{false, true},
				}),
			},
			args: args{
				query: "ESS",
			},
			want: true,
		},
		{
			name: "4",
			fields: fields{
				maze: maze.New([][]bool{
					{true, false},
					{false, true},
				}),
			},
			args: args{
				query: "EWWWSSSSSSNNNSS",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Realization{
				maze: tt.fields.maze,
			}

			got, _ := r.Include(tt.args.query)
			require.Equal(t, tt.want, got)
		})
	}
}
