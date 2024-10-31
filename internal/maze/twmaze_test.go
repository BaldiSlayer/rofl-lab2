package maze

import (
	"fmt"
	"github.com/BaldiSlayer/rofl-lab2/internal/automata"
	"github.com/BaldiSlayer/rofl-lab2/internal/defaults"
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

func TestThinWalled_ToDFA(t *testing.T) {
	type fields struct {
		Maze   [][]LightWallCell
		height int
		width  int
	}
	tests := []struct {
		name   string
		fields fields
		want   *automata.DFA
	}{
		{
			name: "1x1",
			fields: fields{
				width:  1,
				height: 1,
				Maze: [][]LightWallCell{
					{
						{
							leftState:  true,
							rightState: true,
							upState:    true,
							downState:  true,
						},
					},
				},
			},
			want: automata.NewDFA(
				models.Cell{X: 0, Y: 0},
				map[models.Cell]struct{}{
					automata.SpecialState(): {},

					models.Cell{X: -1, Y: -1}: {},
					models.Cell{X: 0, Y: -1}:  {},
					models.Cell{X: 1, Y: -1}:  {},

					models.Cell{X: -1, Y: 0}: {},
					models.Cell{X: 1, Y: 0}:  {},

					models.Cell{X: -1, Y: 1}: {},
					models.Cell{X: 0, Y: 1}:  {},
					models.Cell{X: 1, Y: 1}:  {},
				},
				defaults.GetAlphabet(),
				automata.NewTransitionsFromMap(
					map[models.Cell]map[byte]models.Cell{
						models.Cell{X: -1, Y: -1}: {
							'N': automata.SpecialState(),
							'S': {X: -1, Y: 0},
							'W': automata.SpecialState(),
							'E': {X: 0, Y: -1},
						},
						models.Cell{X: 0, Y: -1}: {
							'N': automata.SpecialState(),
							'S': {X: 0, Y: 0},
							'W': {X: -1, Y: -1},
							'E': {X: 1, Y: -1},
						},
						models.Cell{X: 1, Y: -1}: {
							'N': automata.SpecialState(),
							'S': {X: 1, Y: 0},
							'W': {X: 0, Y: -1},
							'E': automata.SpecialState(),
						},

						models.Cell{X: -1, Y: 0}: {
							'N': {X: -1, Y: -1},
							'S': {X: -1, Y: 1},
							'W': automata.SpecialState(),
							'E': {X: 0, Y: 0},
						},
						models.Cell{X: 0, Y: 0}: {
							'N': {X: 0, Y: -1},
							'S': {X: 0, Y: 1},
							'W': {X: -1, Y: 0},
							'E': {X: 1, Y: 0},
						},
						models.Cell{X: 1, Y: 0}: {
							'N': {X: 1, Y: -1},
							'S': {X: 1, Y: 1},
							'W': {X: 0, Y: 0},
							'E': automata.SpecialState(),
						},

						models.Cell{X: -1, Y: 1}: {
							'N': {X: -1, Y: 0},
							'S': automata.SpecialState(),
							'W': automata.SpecialState(),
							'E': {X: 0, Y: 1},
						},
						models.Cell{X: 0, Y: 1}: {
							'N': {X: 0, Y: 0},
							'S': automata.SpecialState(),
							'W': {X: -1, Y: 1},
							'E': {X: 1, Y: 1},
						},
						models.Cell{X: 1, Y: 1}: {
							'N': {X: 1, Y: 0},
							'S': automata.SpecialState(),
							'W': {X: 0, Y: 1},
							'E': automata.SpecialState(),
						},
					},
				),
				[]models.Cell{
					automata.SpecialState(),

					{X: -1, Y: -1},
					{X: 0, Y: -1},
					{X: 1, Y: -1},

					{X: -1, Y: 0},
					{X: 0, Y: 0},
					{X: 1, Y: 0},

					{X: -1, Y: 1},
					{X: 0, Y: 1},
					{X: 1, Y: 1},
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &ThinWalled{
				Maze:   tt.fields.Maze,
				height: tt.fields.height,
				width:  tt.fields.width,
			}

			got := w.ToDFA()

			require.Equal(t, tt.want, got)
		})
	}
}
