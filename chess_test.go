package chess

import (
	"testing"
)

func TestPieceString(t *testing.T) {
	tests := []struct {
		piece Piece
		want  string
	}{
		{Piece(White | Pawn), "P"},
		{Piece(White | Knight), "N"},
		{Piece(White | Bishop), "B"},
		{Piece(White | Rook), "R"},
		{Piece(White | Queen), "Q"},
		{Piece(White | King), "K"},

		{Piece(Black | Pawn), "p"},
		{Piece(Black | Knight), "n"},
		{Piece(Black | Bishop), "b"},
		{Piece(Black | Rook), "r"},
		{Piece(Black | Queen), "q"},
		{Piece(Black | King), "k"},
	}

	for _, test := range tests {
		if got := test.piece.String(); got != test.want {
			t.Errorf("Piece{Type: %v, Color: %v}.String() = %q, want %q", test.piece.Type(), test.piece.Color(), got, test.want)
		}
	}
}

func TestCoordIsValid(t *testing.T) {
	for nc := 0; nc <= 1; nc++ {
		for f := 0; f < 8; f++ {
			for r := 0; r < 16; r++ {
				coord := Coord((f << 4) | (r << 1) | nc)
				want := nc == 1 && f < 8

				if got := coord.IsValid(); got != want {
					t.Errorf("Coord{%08b}.IsValid() = %v", uint8(coord), got)
				}
			}
		}
	}
}

func TestBoardGet(t *testing.T) {
	tests := []struct {
		coord Coord
		want  Piece
	}{
		{NewCoord("a1"), Piece(White | Rook)},
		{NewCoord("b1"), Piece(White | Knight)},
		{NewCoord("c1"), Piece(White | Bishop)},
		{NewCoord("d1"), Piece(White | Queen)},
		{NewCoord("e1"), Piece(White | King)},
		{NewCoord("f1"), Piece(White | Bishop)},
		{NewCoord("g1"), Piece(White | Knight)},
		{NewCoord("h1"), Piece(White | Rook)},

		{NewCoord("a8"), Piece(Black | Rook)},
		{NewCoord("b8"), Piece(Black | Knight)},
		{NewCoord("c8"), Piece(Black | Bishop)},
		{NewCoord("d8"), Piece(Black | Queen)},
		{NewCoord("e8"), Piece(Black | King)},
		{NewCoord("f8"), Piece(Black | Bishop)},
		{NewCoord("g8"), Piece(Black | Knight)},
		{NewCoord("h8"), Piece(Black | Rook)},

		{NewCoord("b2"), Piece(White | Pawn)},
		{NewCoord("c3"), NoPiece},
		{NewCoord("d4"), NoPiece},
		{NewCoord("e5"), NoPiece},
		{NewCoord("f6"), NoPiece},
		{NewCoord("g7"), Piece(Black | Pawn)},

		{NewCoord("F3"), NoPiece},
		{NewCoord("a0"), NoPiece},
		{NewCoord("h9"), NoPiece},
		{NewCoord("1-1"), NoPiece},
	}

	board := NewBoard()
	for _, test := range tests {
		if got := board.Get(test.coord); got != test.want {
			t.Errorf("Board.Get(Coord{%d, %d}) = %v, want %v", test.coord.File(), test.coord.Rank(), got, test.want)
		}
	}
}

func TestBoardSet(t *testing.T) {
	tests := []struct {
		coord     Coord
		realIndex int
		want      Piece
	}{
		{NewCoord("a1"), 0, Piece(White | Rook)},
		{NewCoord("b1"), 1, Piece(White | Knight)},
		{NewCoord("c1"), 2, Piece(White | Bishop)},
		{NewCoord("d1"), 3, Piece(White | Queen)},
		{NewCoord("e1"), 4, Piece(White | King)},
		{NewCoord("f1"), 5, Piece(White | Bishop)},
		{NewCoord("g1"), 6, Piece(White | Knight)},
		{NewCoord("h1"), 7, Piece(White | Rook)},

		{NewCoord("a8"), 56, Piece(Black | Rook)},
		{NewCoord("b8"), 57, Piece(Black | Knight)},
		{NewCoord("c8"), 58, Piece(Black | Bishop)},
		{NewCoord("d8"), 59, Piece(Black | Queen)},
		{NewCoord("e8"), 60, Piece(Black | King)},
		{NewCoord("f8"), 61, Piece(Black | Bishop)},
		{NewCoord("g8"), 62, Piece(Black | Knight)},
		{NewCoord("h8"), 63, Piece(Black | Rook)},

		{NewCoord("b2"), 9, Piece(White | Pawn)},
		{NewCoord("c3"), 18, NoPiece},
		{NewCoord("d4"), 27, NoPiece},
		{NewCoord("e5"), 36, NoPiece},
		{NewCoord("f6"), 45, NoPiece},
		{NewCoord("g7"), 54, Piece(Black | Pawn)},

		{NewCoord("F3"), 21, NoPiece},
		{NewCoord("a0"), 21, NoPiece},
		{NewCoord("h9"), 21, NoPiece},
		{NewCoord("1-1"), 21, NoPiece},
	}

	board := EmptyBoard()
	for _, test := range tests {
		board.Set(test.coord, test.want)
		if got := board.squares[test.realIndex]; got != test.want {
			t.Errorf("Board.Get(%q) = %v, want %v", test.coord, got, test.want)
		}
	}
}

func TestBoardMakeMove(t *testing.T) {
	sequence := []struct {
		start, target Coord
		legal         bool
	}{
		{NewCoord("d2"), NewCoord("d4"), true},
		{NewCoord("d7"), NewCoord("d5"), true},
		{NewCoord("e2"), NewCoord("e4"), true},
		{NewCoord("d5"), NewCoord("e4"), true},
		{NewCoord("c1"), NewCoord("f4"), true},
		{NewCoord("g8"), NewCoord("f6"), true},
		{NewCoord("b1"), NewCoord("d2"), true},
		{NewCoord("e8"), NewCoord("d7"), true},
		{NewCoord("e1"), NewCoord("d1"), false},

		{NewCoord("c6"), NewCoord("f4"), false}, // move nil to non-nil
		{NewCoord("c5"), NewCoord("c4"), false}, // move nil to nil
		{NewCoord("a0"), NewCoord("f4"), false}, // invalid start
		{NewCoord("a1"), NewCoord("f0"), false}, // invalid end
	}

	board := NewBoard()
	for _, move := range sequence {
		if got := board.MakeMove(Move{move.start, move.target}); got != move.legal {
			t.Errorf("Move{Start: %q, Target: %q} legal: %v, want %v", move.start, move.target, got, move.legal)
		}
	}

	want := Board{
		squares: [64]Piece{
			Piece(White | Rook), Piece(White | Knight), NoPiece, Piece(White | Queen), Piece(White | King), Piece(White | Bishop), Piece(White | Knight), Piece(White | Rook),
			Piece(White | Pawn), Piece(White | Pawn), Piece(White | Pawn), NoPiece, NoPiece, Piece(White | Pawn), Piece(White | Pawn), Piece(White | Pawn),
			NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
			NoPiece, NoPiece, NoPiece, Piece(White | Pawn), Piece(Black | Pawn), Piece(White | Bishop), NoPiece, NoPiece,
			NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
			NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
			Piece(Black | Pawn), Piece(Black | Pawn), Piece(Black | Pawn), NoPiece, Piece(Black | Pawn), Piece(Black | Pawn), Piece(Black | Pawn), Piece(Black | Pawn),
			Piece(Black | Rook), Piece(Black | Knight), Piece(Black | Bishop), Piece(Black | Queen), Piece(Black | King), Piece(Black | Bishop), Piece(Black | Knight), Piece(Black | Rook),
		},
	}

	for i := 0; i < len(board.squares); i++ {
		if board.squares[i] != want.squares[i] {
			t.Log(board)
			t.Errorf("board[%d] = %v, want %v", i, board.squares[i], want.squares[i])
			break
		}
	}
}

func TestBoardToFEN(t *testing.T) {
	tests := []struct {
		board *Board
		want  FEN
	}{
		{
			NewBoard(),
			FEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"),
		},
		{
			&Board{
				squares: [64]Piece{
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, Piece(White | King), NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, Piece(White | Pawn), Piece(White | Pawn), NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, Piece(Black | King), NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, Piece(Black | Knight), NoPiece, NoPiece, Piece(Black | Pawn),
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
				},
				CanCastle:       [4]bool{false, false, false, false},
				FullmoveCounter: 1,
			},
			FEN("8/4n2p/4k3/1PP5/2K5/8/8/8 w - - 0 1"),
		},
		{
			&Board{
				squares: [64]Piece{
					Piece(White | Rook), NoPiece, NoPiece, NoPiece, Piece(White | King), NoPiece, NoPiece, Piece(White | Rook),
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					Piece(Black | Rook), NoPiece, NoPiece, NoPiece, Piece(Black | King), NoPiece, NoPiece, NoPiece,
				},
				CanCastle:       [4]bool{true, true, false, true},
				FullmoveCounter: 1,
			},
			FEN("r3k3/8/8/8/8/8/8/R3K2R w KQq - 0 1"),
		},
		{
			&Board{
				squares: [64]Piece{
					Piece(White | Rook), NoPiece, NoPiece, NoPiece, Piece(White | King), NoPiece, NoPiece, Piece(White | Rook),
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, Piece(Black | King), NoPiece, NoPiece, NoPiece,
				},
				CanCastle:       [4]bool{true, true, false, false},
				FullmoveCounter: 1,
			},
			FEN("4k3/8/8/8/8/8/8/R3K2R w KQ - 0 1"),
		},
		{
			&Board{
				squares: [64]Piece{
					Piece(White | Rook), NoPiece, NoPiece, NoPiece, Piece(White | King), NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, Piece(Black | King), NoPiece, NoPiece, NoPiece,
				},
				SideToMove:      Black,
				CanCastle:       [4]bool{false, true, false, false},
				FullmoveCounter: 1,
			},
			FEN("4k3/8/8/8/8/8/8/R3K3 b Q - 0 1"),
		},
		{
			&Board{
				squares: [64]Piece{
					NoPiece, NoPiece, NoPiece, NoPiece, Piece(White | King), NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, Piece(Black | King), NoPiece, NoPiece, NoPiece,
				},
				SideToMove:      Black,
				CanCastle:       [4]bool{false, false, false, false},
				FullmoveCounter: 1,
			},
			FEN("4k3/8/8/8/8/8/8/4K3 b - - 0 1"),
		},
	}

	for _, test := range tests {
		if got, ok := test.board.ToFEN(); !ok || got != test.want {
			t.Errorf("Board.ToFEN() = %q, want %q", got, test.want)
		}
	}
}

func TestBoardString(t *testing.T) {
	board := NewBoard()
	want := `r n b q k b n r  8
p p p p p p p p  7
. . . . . . . .  6
. . . . . . . .  5
. . . . . . . .  4
. . . . . . . .  3
P P P P P P P P  2
R N B Q K B N R  1
a b c d e f g h`

	board.MakeMove(Move{NewCoord("e2"), NewCoord("e4")})
	board.MakeMove(Move{NewCoord("d7"), NewCoord("d5")})
	board.Flip()
	if got := board.String(); got != want {
		t.Errorf("Board.String() =\n%v, want\n%v", got, want)
	}
}

func TestFENToBoard(t *testing.T) {
	tests := []struct {
		want *Board
		fen  FEN
	}{
		{
			NewBoard(),
			FEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"),
		},
		{
			&Board{
				squares: [64]Piece{
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, Piece(White | King), NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, Piece(White | Pawn), Piece(White | Pawn), NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, Piece(Black | King), NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, Piece(Black | Knight), NoPiece, NoPiece, Piece(Black | Pawn),
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
				},
				CanCastle:       [4]bool{false, false, false, false},
				HalfmoveClock:   5,
				FullmoveCounter: 13,
			},
			FEN("8/4n2p/4k3/1PP5/2K5/8/8/8 w - - 5 13"),
		},
		{
			&Board{
				squares: [64]Piece{
					Piece(White | Rook), NoPiece, NoPiece, NoPiece, Piece(White | King), NoPiece, NoPiece, Piece(White | Rook),
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					Piece(Black | Rook), NoPiece, NoPiece, NoPiece, Piece(Black | King), NoPiece, NoPiece, NoPiece,
				},
				CanCastle:       [4]bool{true, true, false, true},
				FullmoveCounter: 1,
			},
			FEN("r3k3/8/8/8/8/8/8/R3K2R w KQq - 0 1"),
		},
		{
			&Board{
				squares: [64]Piece{
					Piece(White | Rook), NoPiece, NoPiece, NoPiece, Piece(White | King), NoPiece, NoPiece, Piece(White | Rook),
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, Piece(Black | King), NoPiece, NoPiece, NoPiece,
				},
				CanCastle:       [4]bool{true, true, false, false},
				FullmoveCounter: 1,
			},
			FEN("4k3/8/8/8/8/8/8/R3K2R w KQ - 0 1"),
		},
		{
			&Board{
				squares: [64]Piece{
					Piece(White | Rook), NoPiece, NoPiece, NoPiece, Piece(White | King), NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, Piece(Black | King), NoPiece, NoPiece, NoPiece,
				},
				SideToMove:      Black,
				CanCastle:       [4]bool{false, true, false, false},
				FullmoveCounter: 1,
			},
			FEN("4k3/8/8/8/8/8/8/R3K3 b Q - 0 1"),
		},
		{
			&Board{
				squares: [64]Piece{
					NoPiece, NoPiece, NoPiece, NoPiece, Piece(White | King), NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, Piece(Black | King), NoPiece, NoPiece, NoPiece,
				},
				SideToMove:      Black,
				CanCastle:       [4]bool{false, false, false, false},
				FullmoveCounter: 1,
			},
			FEN("4k3/8/8/8/8/8/8/4K3 b - - 0 1"),
		},
		{
			&Board{
				squares: [64]Piece{
					Piece(White | Rook), Piece(White | Knight), Piece(White | Bishop), Piece(White | Queen), Piece(White | King), Piece(White | Bishop), Piece(White | Knight), Piece(White | Rook),
					Piece(White | Pawn), Piece(White | Pawn), Piece(White | Pawn), NoPiece, Piece(White | Pawn), Piece(White | Pawn), Piece(White | Pawn), Piece(White | Pawn),
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, Piece(White | Pawn), Piece(Black | Pawn), NoPiece, NoPiece, NoPiece,
					NoPiece, NoPiece, NoPiece, NoPiece, NoPiece, Piece(Black | Pawn), NoPiece, NoPiece,
					Piece(Black | Pawn), Piece(Black | Pawn), Piece(Black | Pawn), Piece(Black | Pawn), NoPiece, NoPiece, Piece(Black | Pawn), Piece(Black | Pawn),
					Piece(Black | Rook), Piece(Black | Knight), Piece(Black | Bishop), Piece(Black | Queen), Piece(Black | King), Piece(Black | Bishop), Piece(Black | Knight), Piece(Black | Rook),
				},
				CanCastle:       [4]bool{true, true, true, true},
				EnPassantTarget: NewCoord("e6"),
				FullmoveCounter: 3,
			},
			FEN("rnbqkbnr/pppp2pp/5p2/3Pp3/8/8/PPP1PPPP/RNBQKBNR w KQkq e6 0 3"),
		},
		{
			nil,
			FEN("1p a b c d e"),
		},
		{
			nil,
			FEN("1 2 3 4 5"),
		},
	}

	for _, test := range tests {
		got, err := test.fen.ToBoard()

		if got == nil || test.want == nil {
			if got != test.want {
				if err != nil {
					t.Log(err)
				}
				t.Errorf("FEN.ToBoard(%q) =\n%v, want\n%v", test.fen, got, test.want)
			}
		} else {
			// check piece placement
			for i := 0; i < len(got.squares); i++ {
				if got.squares[i] != test.want.squares[i] {
					t.Errorf("FEN.ToBoard(%q) =\n%v, want\n%v", test.fen, got, test.want)
					break
				}
			}

			// check metadata
			if got.SideToMove != test.want.SideToMove {
				t.Errorf("FEN.ToBoard(%q).SideToMove = %v, want %v", test.fen, got.SideToMove, test.want.SideToMove)
			} else if got.CanCastle != test.want.CanCastle {
				t.Errorf("FEN.ToBoard(%q).CanCastle = %v, want %v", test.fen, got.CanCastle, test.want.CanCastle)
			} else if got.EnPassantTarget != test.want.EnPassantTarget {
				t.Errorf("FEN.ToBoard(%q).EnPassantTarget = %v, want %v", test.fen, got.EnPassantTarget, test.want.EnPassantTarget)
			} else if got.HalfmoveClock != test.want.HalfmoveClock {
				t.Errorf("FEN.ToBoard(%q).HalfmoveClock = %v, want %v", test.fen, got.HalfmoveClock, test.want.HalfmoveClock)
			} else if got.FullmoveCounter != test.want.FullmoveCounter {
				t.Errorf("FEN.ToBoard(%q).SideToMove = %v, want %v", test.fen, got.FullmoveCounter, test.want.FullmoveCounter)
			}
		}
	}
}
