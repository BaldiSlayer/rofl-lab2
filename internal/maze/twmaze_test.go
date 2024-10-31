package maze

import (
	"fmt"
	"github.com/BaldiSlayer/rofl-lab2/pkg/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestThinWalled_GetPath(t *testing.T) {
	type fields struct {
		Maze [][]LightWallCell
	}

	type args struct {
		start models.Cell
		end   models.Cell
	}

	maze3x3 := [][]LightWallCell{
		{
			{
				leftState:  false,
				rightState: false,
				upState:    false,
				downState:  true,
			},
			{
				leftState:  false,
				rightState: true,
				upState:    true,
				downState:  true,
			},
			{
				leftState:  true,
				rightState: true,
				upState:    true,
				downState:  true,
			},
		},
		{
			{
				leftState:  false,
				rightState: true,
				upState:    true,
				downState:  false,
			},
			{
				leftState:  true,
				rightState: true,
				upState:    false,
				downState:  true,
			},
			{
				leftState:  true,
				rightState: false,
				upState:    true,
				downState:  false,
			},
		},
		{
			{
				leftState:  true,
				rightState: true,
				upState:    false,
				downState:  true,
			},
			{
				leftState:  true,
				rightState: true,
				upState:    true,
				downState:  false,
			},
			{
				leftState:  true,
				rightState: true,
				upState:    false,
				downState:  true,
			},
		},
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "1",
			fields: fields{
				Maze: maze3x3,
			},
			args: args{
				start: models.Cell{X: 0, Y: 0},
				end:   models.Cell{X: 1, Y: 1},
			},
			want: "SE",
		},
		{
			name: "2",
			fields: fields{
				Maze: maze3x3,
			},
			args: args{
				start: models.Cell{X: 0, Y: 0},
				end:   models.Cell{X: 1, Y: 0},
			},
			want: "SEENW",
		},
		{
			name: "3",
			fields: fields{
				Maze: maze3x3,
			},
			args: args{
				start: models.Cell{X: 0, Y: 0},
				end:   models.Cell{X: 0, Y: 2},
			},
			want: "SESW",
		},
		{
			name: "4",
			fields: fields{
				Maze: maze3x3,
			},
			args: args{
				start: models.Cell{X: 0, Y: 0},
				end:   models.Cell{X: 2, Y: 2},
			},
			want: "SESE",
		},
		{
			name: "5",
			fields: fields{
				Maze: maze3x3,
			},
			args: args{
				start: models.Cell{X: 0, Y: 0},
				end:   models.Cell{X: 0, Y: -1},
			},
			want: "SEENNWW",
		},
		{
			name: "5",
			fields: fields{
				Maze: maze3x3,
			},
			args: args{
				start: models.Cell{X: 0, Y: 0},
				end:   models.Cell{X: -1, Y: -1},
			},
			want: "SESWWNNN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &ThinWalled{
				Maze: tt.fields.Maze,
			}

			got := w.GetPath(tt.args.start, tt.args.end)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestThinWalled_mazeIterator(t *testing.T) {
	r := ThinWalled{height: 2, width: 2}
	visited := make(map[string]bool)

	f := func(y, x int) {
		visited[fmt.Sprintf("%d,%d", y, x)] = true
	}

	r.mazeIterator(bitMaskStorerStore(outCells), f)
	require.Equal(t, 12, len(visited))

	visited = make(map[string]bool)

	r.mazeIterator(bitMaskStorerStore(inCells), f)
	require.Equal(t, 4, len(visited))

	visited = make(map[string]bool)

	r.mazeIterator(bitMaskStorerStore(inCells, outCells), f)
	require.Equal(t, 16, len(visited))
}
