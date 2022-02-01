package chess

import (
	"testing"
)

func TestPieceString(t *testing.T) {
	tests := []struct {
		piece Piece
		want  string
	}{
		{Piece{White, Pawn}, "P"},
		{Piece{White, Knight}, "N"},
		{Piece{White, Bishop}, "B"},
		{Piece{White, Rook}, "R"},
		{Piece{White, Queen}, "Q"},
		{Piece{White, King}, "K"},

		{Piece{Black, Pawn}, "p"},
		{Piece{Black, Knight}, "n"},
		{Piece{Black, Bishop}, "b"},
		{Piece{Black, Rook}, "r"},
		{Piece{Black, Queen}, "q"},
		{Piece{Black, King}, "k"},
	}

	for _, test := range tests {
		if got := test.piece.String(); got != test.want {
			t.Errorf("Piece{Type: %v, Color: %v}.String() = %q, want %q", test.piece.Name, test.piece.Color, got, test.want)
		}
	}
}

func TestBoardAt(t *testing.T) {
	tests := []struct {
		coord Coord
		want  Piece
	}{
		{Coord{1, 1}, Piece{White, Rook}},
		{Coord{2, 1}, Piece{White, Knight}},
		{Coord{3, 1}, Piece{White, Bishop}},
		{Coord{4, 1}, Piece{White, Queen}},
		{Coord{5, 1}, Piece{White, King}},
		{Coord{6, 1}, Piece{White, Bishop}},
		{Coord{7, 1}, Piece{White, Knight}},
		{Coord{8, 1}, Piece{White, Rook}},

		{Coord{1, 8}, Piece{Black, Rook}},
		{Coord{2, 8}, Piece{Black, Knight}},
		{Coord{3, 8}, Piece{Black, Bishop}},
		{Coord{4, 8}, Piece{Black, Queen}},
		{Coord{5, 8}, Piece{Black, King}},
		{Coord{6, 8}, Piece{Black, Bishop}},
		{Coord{7, 8}, Piece{Black, Knight}},
		{Coord{8, 8}, Piece{Black, Rook}},

		{Coord{2, 2}, Piece{White, Pawn}},
		{Coord{3, 3}, Piece{0, 0}},
		{Coord{4, 4}, Piece{0, 0}},
		{Coord{5, 5}, Piece{0, 0}},
		{Coord{6, 6}, Piece{0, 0}},
		{Coord{7, 7}, Piece{Black, Pawn}},
	}

	board := StartingPosition()
	for _, test := range tests {
		if got := board.At(test.coord); *got != test.want {
			t.Errorf("Board.Get(Coord{%d, %d}) = %v, want %v", test.coord.File, test.coord.Rank, got, test.want)
		}
	}
}

// func TestBoardMakeMove(t *testing.T) {
// 	sequence := []struct {
// 		start, target Coord
// 		legal         bool
// 	}{
// 		{CoordFromString("d2"), CoordFromString("d4"), true},
// 		{CoordFromString("d7"), CoordFromString("d5"), true},
// 		{CoordFromString("e2"), CoordFromString("e4"), true},
// 		{CoordFromString("d5"), CoordFromString("e4"), true},
// 		{CoordFromString("c1"), CoordFromString("f4"), true},
// 		{CoordFromString("g8"), CoordFromString("f6"), true},
// 		{CoordFromString("b1"), CoordFromString("d2"), true},
// 		{CoordFromString("e8"), CoordFromString("d7"), true},
// 		{CoordFromString("e1"), CoordFromString("d1"), false},

// 		{CoordFromString("c6"), CoordFromString("f4"), false}, // move nil to non-nil
// 		{CoordFromString("c5"), CoordFromString("c4"), false}, // move nil to nil
// 		{CoordFromString("a0"), CoordFromString("f4"), false}, // invalid start
// 		{CoordFromString("a1"), CoordFromString("f0"), false}, // invalid end
// 	}

// 	board := NewBoard()
// 	for _, move := range sequence {
// 		if got := board.MakeMove(Move{move.start, move.target}); got != move.legal {
// 			t.Errorf("Move{Start: %q, Target: %q} legal: %v, want %v", move.start, move.target, got, move.legal)
// 		}
// 	}

// 	want := Board{
// 		squares: [64]Piece{
// 			Piece{White, Rook}, Piece{White, Knight}, Piece{0, 0}, Piece{White, Queen}, Piece{White, King}, Piece{White, Bishop}, Piece{White, Knight}, Piece{White, Rook},
// 			Piece{White, Pawn}, Piece{White, Pawn}, Piece{White, Pawn}, Piece{0, 0}, Piece{0, 0}, Piece{White, Pawn}, Piece{White, Pawn}, Piece{White, Pawn},
// 			Piece{0, 0}, Piece{0, 0}, Piece{0, 0}, Piece{0, 0}, Piece{0, 0}, Piece{0, 0}, Piece{0, 0}, Piece{0, 0},
// 			Piece{0, 0}, Piece{0, 0}, Piece{0, 0}, Piece{White, Pawn}, Piece{Black, Pawn}, Piece{White, Bishop}, Piece{0, 0}, Piece{0, 0},
// 			Piece{0, 0}, Piece{0, 0}, Piece{0, 0}, Piece{0, 0}, Piece{0, 0}, Piece{0, 0}, Piece{0, 0}, Piece{0, 0},
// 			Piece{0, 0}, Piece{0, 0}, Piece{0, 0}, Piece{0, 0}, Piece{0, 0}, Piece{0, 0}, Piece{0, 0}, Piece{0, 0},
// 			Piece{Black, Pawn}, Piece{Black, Pawn}, Piece{Black, Pawn}, Piece{0, 0}, Piece{Black, Pawn}, Piece{Black, Pawn}, Piece{Black, Pawn}, Piece{Black, Pawn},
// 			Piece{Black, Rook}, Piece{Black, Knight}, Piece{Black, Bishop}, Piece{Black, Queen}, Piece{Black, King}, Piece{Black, Bishop}, Piece{Black, Knight}, Piece{Black, Rook},
// 		},
// 	}

// 	for i := 0; i < len(board.squares); i++ {
// 		if board.squares[i] != want.squares[i] {
// 			t.Log(board)
// 			t.Errorf("board[%d] = %v, want %v", i, board.squares[i], want.squares[i])
// 			break
// 		}
// 	}
// }

func TestCastles(t *testing.T) {
	castles := Castles(0)

	if castles&15 != 0 {
		t.Errorf("Castles.Can(White, Kingside) != 0")
	}

	castles.Allow(Black, Kingside)
	if castles&15 != 0b0100 {
		t.Errorf("Castles.Allow(Black, Kingside) = %04b, want 0100", castles)
	}

	if got := castles.Can(Black, Kingside); got == false {
		t.Error("Castles.Can(Black, Kingside) != true")
	}

	castles.Disallow(Black, Kingside)
	if castles&15 != 0 {
		t.Errorf("Castles.Disallow(Black, Kingside) != 0")
	}
}

func TestBoardToFEN(t *testing.T) {
	tests := []struct {
		board *Board
		want  string
	}{
		{
			StartingPosition(),
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		},
		{
			&Board{
				squares: [64]Piece{
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {White, King}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {White, Pawn}, {White, Pawn}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {Black, King}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {Black, Knight}, {0, 0}, {0, 0}, {Black, Pawn},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
				},
				SideToMove:      White,
				CastleRights:    NewCastles("-"),
				FullmoveCounter: 1,
			},
			"8/4n2p/4k3/1PP5/2K5/8/8/8 w - - 0 1",
		},
		{
			&Board{
				squares: [64]Piece{
					{White, Rook}, {0, 0}, {0, 0}, {0, 0}, {White, King}, {0, 0}, {0, 0}, {White, Rook},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{Black, Rook}, {0, 0}, {0, 0}, {0, 0}, {Black, King}, {0, 0}, {0, 0}, {0, 0},
				},
				SideToMove:      White,
				CastleRights:    NewCastles("KQq"),
				FullmoveCounter: 1,
			},
			"r3k3/8/8/8/8/8/8/R3K2R w KQq - 0 1",
		},
		{
			&Board{
				squares: [64]Piece{
					{White, Rook}, {0, 0}, {0, 0}, {0, 0}, {White, King}, {0, 0}, {0, 0}, {White, Rook},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {Black, King}, {0, 0}, {0, 0}, {0, 0},
				},
				SideToMove:      White,
				CastleRights:    NewCastles("KQ"),
				FullmoveCounter: 1,
			},
			"4k3/8/8/8/8/8/8/R3K2R w KQ - 0 1",
		},
		{
			&Board{
				squares: [64]Piece{
					{White, Rook}, {0, 0}, {0, 0}, {0, 0}, {White, King}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {Black, King}, {0, 0}, {0, 0}, {0, 0},
				},
				SideToMove:      Black,
				CastleRights:    NewCastles("Q"),
				FullmoveCounter: 1,
			},
			"4k3/8/8/8/8/8/8/R3K3 b Q - 0 1",
		},
		{
			&Board{
				squares: [64]Piece{
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {White, King}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {Black, King}, {0, 0}, {0, 0}, {0, 0},
				},
				SideToMove:      Black,
				CastleRights:    NewCastles("-"),
				FullmoveCounter: 1,
			},
			"4k3/8/8/8/8/8/8/4K3 b - - 0 1",
		},
	}

	for _, test := range tests {
		if got := test.board.String(); got != test.want {
			t.Errorf("Board.String() = %q, want %q", got, test.want)
		}
	}
}

func TestFENToBoard(t *testing.T) {
	tests := []struct {
		want *Board
		fen  string
	}{
		{
			StartingPosition(),
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		},
		{
			&Board{
				squares: [64]Piece{
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {White, King}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {White, Pawn}, {White, Pawn}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {Black, King}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {Black, Knight}, {0, 0}, {0, 0}, {Black, Pawn},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
				},
				SideToMove:      White,
				CastleRights:    NewCastles("-"),
				HalfmoveClock:   5,
				FullmoveCounter: 13,
			},
			"8/4n2p/4k3/1PP5/2K5/8/8/8 w - - 5 13",
		},
		{
			&Board{
				squares: [64]Piece{
					{White, Rook}, {0, 0}, {0, 0}, {0, 0}, {White, King}, {0, 0}, {0, 0}, {White, Rook},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{Black, Rook}, {0, 0}, {0, 0}, {0, 0}, {Black, King}, {0, 0}, {0, 0}, {0, 0},
				},
				SideToMove:      White,
				CastleRights:    NewCastles("KQq"),
				FullmoveCounter: 1,
			},
			"r3k3/8/8/8/8/8/8/R3K2R w KQq - 0 1",
		},
		{
			&Board{
				squares: [64]Piece{
					{White, Rook}, {0, 0}, {0, 0}, {0, 0}, {White, King}, {0, 0}, {0, 0}, {White, Rook},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {Black, King}, {0, 0}, {0, 0}, {0, 0},
				},
				SideToMove:      White,
				CastleRights:    NewCastles("KQ"),
				FullmoveCounter: 1,
			},
			"4k3/8/8/8/8/8/8/R3K2R w KQ - 0 1",
		},
		{
			&Board{
				squares: [64]Piece{
					{White, Rook}, {0, 0}, {0, 0}, {0, 0}, {White, King}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {Black, King}, {0, 0}, {0, 0}, {0, 0},
				},
				SideToMove:      Black,
				CastleRights:    NewCastles("Q"),
				FullmoveCounter: 1,
			},
			"4k3/8/8/8/8/8/8/R3K3 b Q - 0 1",
		},
		{
			&Board{
				squares: [64]Piece{
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {White, King}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {Black, King}, {0, 0}, {0, 0}, {0, 0},
				},
				SideToMove:      Black,
				CastleRights:    NewCastles("-"),
				FullmoveCounter: 1,
			},
			"4k3/8/8/8/8/8/8/4K3 b - - 0 1",
		},
		{
			&Board{
				squares: [64]Piece{
					{White, Rook}, {White, Knight}, {White, Bishop}, {White, Queen}, {White, King}, {White, Bishop}, {White, Knight}, {White, Rook},
					{White, Pawn}, {White, Pawn}, {White, Pawn}, {0, 0}, {White, Pawn}, {White, Pawn}, {White, Pawn}, {White, Pawn},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {White, Pawn}, {Black, Pawn}, {0, 0}, {0, 0}, {0, 0},
					{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {Black, Pawn}, {0, 0}, {0, 0},
					{Black, Pawn}, {Black, Pawn}, {Black, Pawn}, {Black, Pawn}, {0, 0}, {0, 0}, {Black, Pawn}, {Black, Pawn},
					{Black, Rook}, {Black, Knight}, {Black, Bishop}, {Black, Queen}, {Black, King}, {Black, Bishop}, {Black, Knight}, {Black, Rook},
				},
				SideToMove:      White,
				CastleRights:    NewCastles("KQkq"),
				EnPassantTarget: NewCoord("e6"),
				FullmoveCounter: 3,
			},
			"rnbqkbnr/pppp2pp/5p2/3Pp3/8/8/PPP1PPPP/RNBQKBNR w KQkq e6 0 3",
		},
		{
			nil,
			"1p a b c d e",
		},
		{
			nil,
			"1 2 3 4 5",
		},
	}

	for _, test := range tests {
		got, err := NewBoard(test.fen)

		if got == nil || test.want == nil {
			if got != test.want {
				if err != nil {
					t.Log(err)
				}
				t.Errorf("NewBoard(%q) =\n%v, want\n%v", test.fen, got, test.want)
			}
		} else {
			// check piece placement
			for i := 0; i < len(got.squares); i++ {
				if got.squares[i] != test.want.squares[i] {
					t.Errorf("NewBoard(%q) =\n%v, want\n%v", test.fen, got, test.want)
					break
				}
			}

			// check metadata
			if got.SideToMove != test.want.SideToMove {
				t.Errorf("NewBoard(%q).SideToMove = %v, want %v", test.fen, got.SideToMove, test.want.SideToMove)
			} else if got.CastleRights != test.want.CastleRights {
				t.Errorf("NewBoard(%q).CastleRights = %v, want %v", test.fen, got.CastleRights, test.want.CastleRights)
			} else if got.EnPassantTarget != test.want.EnPassantTarget {
				t.Errorf("NewBoard(%q).EnPassantTarget = %v, want %v", test.fen, got.EnPassantTarget, test.want.EnPassantTarget)
			} else if got.HalfmoveClock != test.want.HalfmoveClock {
				t.Errorf("NewBoard(%q).HalfmoveClock = %v, want %v", test.fen, got.HalfmoveClock, test.want.HalfmoveClock)
			} else if got.FullmoveCounter != test.want.FullmoveCounter {
				t.Errorf("NewBoard(%q).SideToMove = %v, want %v", test.fen, got.FullmoveCounter, test.want.FullmoveCounter)
			}
		}
	}
}

func TestCoord(t *testing.T) {
	tests := []struct {
		coord     Coord
		wantFile  int
		wantRank  int
		wantIndex int
		valid     bool
	}{
		{Coord{1, 1}, 1, 1, 0, true},
		{Coord{2, 1}, 2, 1, 1, true},
		{Coord{3, 2}, 3, 2, 10, true},
		{Coord{4, 2}, 4, 2, 11, true},
		{Coord{1, 3}, 1, 3, 16, true},
		{Coord{1, 4}, 1, 4, 24, true},
		{Coord{8, 5}, 8, 5, 39, true},
		{Coord{-1, 9}, 0, 0, 0, false},
	}

	for _, test := range tests {
		if test.coord.IsValid() != test.valid {
			t.Errorf("Coord{File: %d, Rank: %d}.IsValid() != %v", test.wantFile, test.wantRank, test.valid)
		}
		gotFile := test.coord.File
		gotRank := test.coord.Rank
		gotIndex := test.coord.Index()
		if test.coord.IsValid() && (gotFile != test.wantFile || gotRank != test.wantRank || gotIndex != test.wantIndex) {
			t.Errorf("Coord{File: %d, Rank: %d, Index: %d} != {%d, %d, %d}", gotFile, gotRank, gotIndex, test.wantFile, test.wantRank, test.wantIndex)
		}
	}
}

func TestUnmakeMove(t *testing.T) {
	tests := []struct {
		position string
		move     Move
	}{
		{
			"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			Move{Coord{5, 1}, Coord{7, 1}, MoveFlags{Moves: King, Castle: Kingside}},
		},
	}

	for _, test := range tests {
		board, _ := NewBoard(test.position)
		apply := *board

		if err := apply.MakeMove(test.move); err != nil {
			t.Errorf("Make move failed: %v", err)
		}
		apply.UnmakeMove()

		if board.squares != apply.squares {
			t.Errorf("board squares don't match: %q != %q", board.String(), apply.String())
		}
		if board.SideToMove != apply.SideToMove {
			t.Errorf("moves sides don't match: %v != %v", board.SideToMove, apply.SideToMove)
		}
		if board.CastleRights != apply.CastleRights {
			t.Errorf("castle rights don't match: %v != %v", board.CastleRights, apply.CastleRights)
		}
		if board.EnPassantTarget != apply.EnPassantTarget {
			t.Errorf("en passant targets don't match: %v != %v", board.EnPassantTarget, apply.EnPassantTarget)
		}
		if board.HalfmoveClock != apply.HalfmoveClock {
			t.Errorf("halfmove clocks don't match: %v != %v", board.HalfmoveClock, apply.HalfmoveClock)
		}
		if board.FullmoveCounter != apply.FullmoveCounter {
			t.Errorf("move counters don't match: %v != %v", board.FullmoveCounter, apply.FullmoveCounter)
		}
	}
}

func TestCountMoves(t *testing.T) {
	tests := []struct {
		position string
		depth    int
		want     int
	}{
		{
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			1,
			20,
		},
		{
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			2,
			400,
		},
		{
			"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			3,
			8902,
		},
		{
			"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			1,
			48,
		},
		{
			"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			2,
			2039,
		},
		{
			"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			3,
			97862,
		},
		{
			"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			4,
			4085603,
		},
		{
			"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
			1,
			14,
		},
		{
			"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
			2,
			191,
		},
		{
			"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
			3,
			2812,
		},
		{
			"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
			6,
			11030083,
		},
		{
			"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
			1,
			6,
		},
		{
			"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
			2,
			264,
		},
		{
			"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
			3,
			9467,
		},
		{
			"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
			4,
			422333,
		},
		{
			"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
			5,
			15833292,
		},
	}

	for _, test := range tests {
		board, _ := NewBoard(test.position)
		if got, breakdown := board.CountMoves(test.depth); got != test.want {
			t.Errorf("CountMoves() on [d=%d] %q = %d, want %d\n\t%v", test.depth, test.position, got, test.want, breakdown)
		}
	}
}

func BenchmarkMoveGen(b *testing.B) {
	board, _ := NewBoard("r2qr1k1/pp3pp1/2n2n1p/2bp4/6b1/2PB1NN1/PP3PPP/R1BQR1K1 w - - 3 13")
	for i := 0; i < b.N; i++ {
		board.GenMoves()
	}
}

func BenchmarkMakeUnmake(b *testing.B) {
	board, _ := NewBoard("r2qr1k1/pp3pp1/2n2n1p/2bp4/6b1/2PB1NN1/PP3PPP/R1BQR1K1 w - - 3 13")
	move := Move{Coord{4, 0}, Coord{4, 8}, MoveFlags{Moves: Rook, Captures: Rook}}

	for i := 0; i < b.N; i++ {
		board.MakeMove(move)
		board.UnmakeMove()
	}
}

func BenchmarkCopyUndo(b *testing.B) {
	board, _ := NewBoard("r2qr1k1/pp3pp1/2n2n1p/2bp4/6b1/2PB1NN1/PP3PPP/R1BQR1K1 w - - 3 13")
	move := Move{Coord{4, 0}, Coord{4, 8}, MoveFlags{Moves: Rook, Captures: Rook}}

	for i := 0; i < b.N; i++ {
		temp := *board
		temp.MakeMove(move)
	}
}

func BenchmarkCountMoves(b *testing.B) {
	board := StartingPosition()

	for i := 0; i < b.N; i++ {
		board.CountMoves(1)
	}
}
