package chess

import (
	"testing"
)

func TestPieceString(t *testing.T) {
	tests := []struct {
		piece Piece
		want  string
	}{
		{NewPiece(White, Pawn), "P"},
		{NewPiece(White, Knight), "N"},
		{NewPiece(White, Bishop), "B"},
		{NewPiece(White, Rook), "R"},
		{NewPiece(White, Queen), "Q"},
		{NewPiece(White, King), "K"},

		{NewPiece(Black, Pawn), "p"},
		{NewPiece(Black, Knight), "n"},
		{NewPiece(Black, Bishop), "b"},
		{NewPiece(Black, Rook), "r"},
		{NewPiece(Black, Queen), "q"},
		{NewPiece(Black, King), "k"},
	}

	for _, test := range tests {
		if got := test.piece.String(); got != test.want {
			t.Errorf("Piece{Type: %v, Color: %v}.String() = %q, want %q", test.piece.Type(), test.piece.Color(), got, test.want)
		}
	}
}

func TestBoardAt(t *testing.T) {
	tests := []struct {
		coord Coord
		want  Piece
	}{
		{NewCoord(0, 0), NewPiece(White, Rook)},
		{NewCoord(1, 0), NewPiece(White, Knight)},
		{NewCoord(2, 0), NewPiece(White, Bishop)},
		{NewCoord(3, 0), NewPiece(White, Queen)},
		{NewCoord(4, 0), NewPiece(White, King)},
		{NewCoord(5, 0), NewPiece(White, Bishop)},
		{NewCoord(6, 0), NewPiece(White, Knight)},
		{NewCoord(7, 0), NewPiece(White, Rook)},

		{NewCoord(0, 7), NewPiece(Black, Rook)},
		{NewCoord(1, 7), NewPiece(Black, Knight)},
		{NewCoord(2, 7), NewPiece(Black, Bishop)},
		{NewCoord(3, 7), NewPiece(Black, Queen)},
		{NewCoord(4, 7), NewPiece(Black, King)},
		{NewCoord(5, 7), NewPiece(Black, Bishop)},
		{NewCoord(6, 7), NewPiece(Black, Knight)},
		{NewCoord(7, 7), NewPiece(Black, Rook)},

		{NewCoord(1, 1), NewPiece(White, Pawn)},
		{NewCoord(2, 2), Piece(0)},
		{NewCoord(3, 3), Piece(0)},
		{NewCoord(4, 4), Piece(0)},
		{NewCoord(5, 5), Piece(0)},
		{NewCoord(6, 6), NewPiece(Black, Pawn)},
	}

	board := StartingPosition()
	for _, test := range tests {
		if got := board.At(test.coord); *got != test.want {
			t.Errorf("Board.Get(Coord{%d, %d}) = %v, want %v", test.coord.File(), test.coord.Rank(), got, test.want)
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
// 			NewPiece(White, Rook), NewPiece(White, Knight), Piece(0), NewPiece(White, Queen), NewPiece(White, King), NewPiece(White, Bishop), NewPiece(White, Knight), NewPiece(White, Rook),
// 			NewPiece(White, Pawn), NewPiece(White, Pawn), NewPiece(White, Pawn), Piece(0), Piece(0), NewPiece(White, Pawn), NewPiece(White, Pawn), NewPiece(White, Pawn),
// 			Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
// 			Piece(0), Piece(0), Piece(0), NewPiece(White, Pawn), NewPiece(Black, Pawn), NewPiece(White, Bishop), Piece(0), Piece(0),
// 			Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
// 			Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
// 			NewPiece(Black, Pawn), NewPiece(Black, Pawn), NewPiece(Black, Pawn), Piece(0), NewPiece(Black, Pawn), NewPiece(Black, Pawn), NewPiece(Black, Pawn), NewPiece(Black, Pawn),
// 			NewPiece(Black, Rook), NewPiece(Black, Knight), NewPiece(Black, Bishop), NewPiece(Black, Queen), NewPiece(Black, King), NewPiece(Black, Bishop), NewPiece(Black, Knight), NewPiece(Black, Rook),
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
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), NewPiece(White, King), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), NewPiece(White, Pawn), NewPiece(White, Pawn), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), NewPiece(Black, King), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), NewPiece(Black, Knight), Piece(0), Piece(0), NewPiece(Black, Pawn),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
				},
				SideToMove:      White,
				CastleRights:    [4]bool{false, false, false, false},
				FullmoveCounter: 1,
			},
			"8/4n2p/4k3/1PP5/2K5/8/8/8 w - - 0 1",
		},
		{
			&Board{
				squares: [64]Piece{
					NewPiece(White, Rook), Piece(0), Piece(0), Piece(0), NewPiece(White, King), Piece(0), Piece(0), NewPiece(White, Rook),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					NewPiece(Black, Rook), Piece(0), Piece(0), Piece(0), NewPiece(Black, King), Piece(0), Piece(0), Piece(0),
				},
				SideToMove:      White,
				CastleRights:    [4]bool{true, true, false, true},
				FullmoveCounter: 1,
			},
			"r3k3/8/8/8/8/8/8/R3K2R w KQq - 0 1",
		},
		{
			&Board{
				squares: [64]Piece{
					NewPiece(White, Rook), Piece(0), Piece(0), Piece(0), NewPiece(White, King), Piece(0), Piece(0), NewPiece(White, Rook),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), NewPiece(Black, King), Piece(0), Piece(0), Piece(0),
				},
				SideToMove:      White,
				CastleRights:    [4]bool{true, true, false, false},
				FullmoveCounter: 1,
			},
			"4k3/8/8/8/8/8/8/R3K2R w KQ - 0 1",
		},
		{
			&Board{
				squares: [64]Piece{
					NewPiece(White, Rook), Piece(0), Piece(0), Piece(0), NewPiece(White, King), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), NewPiece(Black, King), Piece(0), Piece(0), Piece(0),
				},
				SideToMove:      Black,
				CastleRights:    [4]bool{false, true, false, false},
				FullmoveCounter: 1,
			},
			"4k3/8/8/8/8/8/8/R3K3 b Q - 0 1",
		},
		{
			&Board{
				squares: [64]Piece{
					Piece(0), Piece(0), Piece(0), Piece(0), NewPiece(White, King), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), NewPiece(Black, King), Piece(0), Piece(0), Piece(0),
				},
				SideToMove:      Black,
				CastleRights:    [4]bool{false, false, false, false},
				FullmoveCounter: 1,
			},
			"4k3/8/8/8/8/8/8/4K3 b - - 0 1",
		},
	}

	for _, test := range tests {
		if got := test.board.FEN(); got != test.want {
			t.Errorf("Board.FEN() = %q, want %q", got, test.want)
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
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), NewPiece(White, King), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), NewPiece(White, Pawn), NewPiece(White, Pawn), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), NewPiece(Black, King), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), NewPiece(Black, Knight), Piece(0), Piece(0), NewPiece(Black, Pawn),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
				},
				SideToMove:      White,
				CastleRights:    [4]bool{false, false, false, false},
				HalfmoveClock:   5,
				FullmoveCounter: 13,
			},
			"8/4n2p/4k3/1PP5/2K5/8/8/8 w - - 5 13",
		},
		{
			&Board{
				squares: [64]Piece{
					NewPiece(White, Rook), Piece(0), Piece(0), Piece(0), NewPiece(White, King), Piece(0), Piece(0), NewPiece(White, Rook),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					NewPiece(Black, Rook), Piece(0), Piece(0), Piece(0), NewPiece(Black, King), Piece(0), Piece(0), Piece(0),
				},
				SideToMove:      White,
				CastleRights:    [4]bool{true, true, false, true},
				FullmoveCounter: 1,
			},
			"r3k3/8/8/8/8/8/8/R3K2R w KQq - 0 1",
		},
		{
			&Board{
				squares: [64]Piece{
					NewPiece(White, Rook), Piece(0), Piece(0), Piece(0), NewPiece(White, King), Piece(0), Piece(0), NewPiece(White, Rook),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), NewPiece(Black, King), Piece(0), Piece(0), Piece(0),
				},
				SideToMove:      White,
				CastleRights:    [4]bool{true, true, false, false},
				FullmoveCounter: 1,
			},
			"4k3/8/8/8/8/8/8/R3K2R w KQ - 0 1",
		},
		{
			&Board{
				squares: [64]Piece{
					NewPiece(White, Rook), Piece(0), Piece(0), Piece(0), NewPiece(White, King), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), NewPiece(Black, King), Piece(0), Piece(0), Piece(0),
				},
				SideToMove:      Black,
				CastleRights:    [4]bool{false, true, false, false},
				FullmoveCounter: 1,
			},
			"4k3/8/8/8/8/8/8/R3K3 b Q - 0 1",
		},
		{
			&Board{
				squares: [64]Piece{
					Piece(0), Piece(0), Piece(0), Piece(0), NewPiece(White, King), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), NewPiece(Black, King), Piece(0), Piece(0), Piece(0),
				},
				SideToMove:      Black,
				CastleRights:    [4]bool{false, false, false, false},
				FullmoveCounter: 1,
			},
			"4k3/8/8/8/8/8/8/4K3 b - - 0 1",
		},
		{
			&Board{
				squares: [64]Piece{
					NewPiece(White, Rook), NewPiece(White, Knight), NewPiece(White, Bishop), NewPiece(White, Queen), NewPiece(White, King), NewPiece(White, Bishop), NewPiece(White, Knight), NewPiece(White, Rook),
					NewPiece(White, Pawn), NewPiece(White, Pawn), NewPiece(White, Pawn), Piece(0), NewPiece(White, Pawn), NewPiece(White, Pawn), NewPiece(White, Pawn), NewPiece(White, Pawn),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), NewPiece(White, Pawn), NewPiece(Black, Pawn), Piece(0), Piece(0), Piece(0),
					Piece(0), Piece(0), Piece(0), Piece(0), Piece(0), NewPiece(Black, Pawn), Piece(0), Piece(0),
					NewPiece(Black, Pawn), NewPiece(Black, Pawn), NewPiece(Black, Pawn), NewPiece(Black, Pawn), Piece(0), Piece(0), NewPiece(Black, Pawn), NewPiece(Black, Pawn),
					NewPiece(Black, Rook), NewPiece(Black, Knight), NewPiece(Black, Bishop), NewPiece(Black, Queen), NewPiece(Black, King), NewPiece(Black, Bishop), NewPiece(Black, Knight), NewPiece(Black, Rook),
				},
				SideToMove:      White,
				CastleRights:    [4]bool{true, true, true, true},
				EnPassantTarget: NewCoord(4, 5),
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
		got, err := BoardFromString(test.fen)

		if got == nil || test.want == nil {
			if got != test.want {
				if err != nil {
					t.Log(err)
				}
				t.Errorf("BoardFromString(%q) =\n%v, want\n%v", test.fen, got, test.want)
			}
		} else {
			// check piece placement
			for i := 0; i < len(got.squares); i++ {
				if got.squares[i] != test.want.squares[i] {
					t.Errorf("BoardFromString(%q) =\n%v, want\n%v", test.fen, got, test.want)
					break
				}
			}

			// check metadata
			if got.SideToMove != test.want.SideToMove {
				t.Errorf("BoardFromString(%q).SideToMove = %v, want %v", test.fen, got.SideToMove, test.want.SideToMove)
			} else if got.CastleRights != test.want.CastleRights {
				t.Errorf("BoardFromString(%q).CastleRights = %v, want %v", test.fen, got.CastleRights, test.want.CastleRights)
			} else if got.EnPassantTarget != test.want.EnPassantTarget {
				t.Errorf("BoardFromString(%q).EnPassantTarget = %v, want %v", test.fen, got.EnPassantTarget, test.want.EnPassantTarget)
			} else if got.HalfmoveClock != test.want.HalfmoveClock {
				t.Errorf("BoardFromString(%q).HalfmoveClock = %v, want %v", test.fen, got.HalfmoveClock, test.want.HalfmoveClock)
			} else if got.FullmoveCounter != test.want.FullmoveCounter {
				t.Errorf("BoardFromString(%q).SideToMove = %v, want %v", test.fen, got.FullmoveCounter, test.want.FullmoveCounter)
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
	}{
		{NewCoord(0, 0), 0, 0, 0},
		{NewCoord(1, 0), 1, 0, 1},
		{NewCoord(2, 1), 2, 1, 10},
		{NewCoord(3, 1), 3, 1, 11},
		{NewCoord(0, 2), 0, 2, 16},
		{NewCoord(0, 3), 0, 3, 24},
		{NewCoord(7, 4), 7, 4, 39},
		{NewCoord(-1, 9), 0, 0, 0},
	}

	for _, test := range tests {
		gotFile := test.coord.File()
		gotRank := test.coord.Rank()
		gotIndex := test.coord.index()
		if gotFile != test.wantFile || gotRank != test.wantRank || gotIndex != test.wantIndex {
			t.Errorf("Coord{File: %d, Rank: %d, Index: %d} != {%d, %d, %d}", gotFile, gotRank, gotIndex, test.wantFile, test.wantRank, test.wantIndex)
		}
	}
}

func BenchmarkMoveGen(b *testing.B) {
	board, _ := BoardFromString("r2qr1k1/pp3pp1/2n2n1p/2bp4/6b1/2PB1NN1/PP3PPP/R1BQR1K1 w - - 3 13")
	for i := 0; i < b.N; i++ {
		board.GenMoves()
	}
}
