package chess

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

/*
 *
 *
 * Piece structure methods / functions
 */
func NewPiece(symbol rune) Piece {
	return map[rune]Piece{
		'p': Piece(Black | Pawn),
		'n': Piece(Black | Knight),
		'b': Piece(Black | Bishop),
		'r': Piece(Black | Rook),
		'q': Piece(Black | Queen),
		'k': Piece(Black | King),

		'P': Piece(White | Pawn),
		'N': Piece(White | Knight),
		'B': Piece(White | Bishop),
		'R': Piece(White | Rook),
		'Q': Piece(White | Queen),
		'K': Piece(White | King),
	}[symbol]
}
func (piece Piece) Type() uint8 {
	return uint8(piece) & typeMask
}
func (piece Piece) Color() uint8 {
	return uint8(piece) & colorMask
}
func (piece Piece) String() string {
	if piece > Piece(Black|King) || piece == NoPiece {
		return ""
	}

	i := (piece.Type() >> typeShift) - 1
	name := [...]rune{'p', 'n', 'b', 'r', 'q', 'k'}[i]
	if (piece.Color() >> colorShift) == White {
		name += 'A' - 'a'
	}
	return string(name)
}

/*
 *
 *
 * Coord structure methods / functions
 */
func NewCoord(pos string) Coord {
	if len(pos) != 2 || pos[0] < 'a' || pos[0] > 'h' || pos[1] < '1' || pos[1] > '8' {
		return NoCoord
	}

	return Coord((uint8(pos[0]-'a') << fileShift) | (uint8(pos[1]-'1') << rankShift) | 1)
}
func (coord Coord) IsValid() bool {
	return uint8(coord)&1 == 1 && uint8(coord) <= File|Rank|1
}
func (coord Coord) File() int {
	return int((uint8(coord) & File) >> fileShift)
}
func (coord Coord) Rank() int {
	return int((uint8(coord) & Rank) >> rankShift)
}
func (coord Coord) String() string {
	return fmt.Sprintf("%c%d", 'a'+rune(coord.File()), coord.Rank()+1)
}

func (move Move) IsValid(board *Board) bool {
	return move.Start.IsValid() && move.Target.IsValid() && board.Get(move.Start) != NoPiece
}
func (move Move) Translation() (x, y int) {
	return move.Target.File() - move.Start.File(), move.Target.Rank() - move.Start.Rank()
}

/*
 *
 *
 * Board structure methods / functions
 */
func NewBoard() *Board {
	return &Board{
		squares: [64]Piece{
			Piece(White | Rook), Piece(White | Knight), Piece(White | Bishop), Piece(White | Queen), Piece(White | King), Piece(White | Bishop), Piece(White | Knight), Piece(White | Rook),
			Piece(White | Pawn), Piece(White | Pawn), Piece(White | Pawn), Piece(White | Pawn), Piece(White | Pawn), Piece(White | Pawn), Piece(White | Pawn), Piece(White | Pawn),
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			Piece(Black | Pawn), Piece(Black | Pawn), Piece(Black | Pawn), Piece(Black | Pawn), Piece(Black | Pawn), Piece(Black | Pawn), Piece(Black | Pawn), Piece(Black | Pawn),
			Piece(Black | Rook), Piece(Black | Knight), Piece(Black | Bishop), Piece(Black | Queen), Piece(Black | King), Piece(Black | Bishop), Piece(Black | Knight), Piece(Black | Rook),
		},
		CanCastle:       [4]bool{true, true, true, true},
		FullmoveCounter: 1,
	}
}
func EmptyBoard() (board *Board) {
	defer func() { board.FullmoveCounter = 1 }()
	return new(Board)
}
func (board *Board) Get(coord Coord) Piece {
	if !coord.IsValid() {
		return NoPiece
	}
	return board.squares[coord.Rank()*8+coord.File()]
}
func (board *Board) Set(coord Coord, piece Piece) {
	if !coord.IsValid() {
		return
	}
	board.squares[coord.Rank()*8+coord.File()] = piece

}
func (board *Board) Flip() {
	if board.Orientation == White {
		board.Orientation = Black
	} else {
		board.Orientation = White
	}
}
func assembleFEN(board *Board, placement string) FEN {
	moveSide := "w"
	if board.SideToMove == Black {
		moveSide = "b"
	}

	castleLetters := make([]string, 4)
	if board.CanCastle[WKingside] {
		castleLetters[0] = "K"
	}
	if board.CanCastle[WQueenside] {
		castleLetters[1] = "Q"
	}
	if board.CanCastle[BKingside] {
		castleLetters[2] = "k"
	}
	if board.CanCastle[BQueenside] {
		castleLetters[3] = "q"
	}
	castlingAbility := strings.Join(castleLetters, "")
	if castlingAbility == "" {
		castlingAbility = "-"
	}

	enPassantTarget := "-"
	if board.EnPassantTarget != NoCoord {
		enPassantTarget = board.EnPassantTarget.String()
	}

	return FEN(fmt.Sprintf(
		"%v %v %v %v %v %v",
		placement,
		moveSide,
		castlingAbility,
		enPassantTarget,
		board.HalfmoveClock,
		board.FullmoveCounter,
	))
}
func (board *Board) ToFEN() (FEN, bool) {
	placement := []string{}

	spaces := 0
	for f, r := 0, 7; r >= 0; {
		piece := board.squares[r*8+f]
		if piece == NoPiece {
			spaces++
		} else {
			if pieceString := piece.String(); pieceString == "" {
				return "", false
			} else {
				if spaces > 0 {
					placement = append(placement, fmt.Sprint(spaces))
				}
				placement = append(placement, pieceString)
				spaces = 0
			}
		}

		f++
		if f&7 == 0 { // file is divisible by 8
			r -= 1
			f = 0

			if spaces > 0 {
				placement = append(placement, fmt.Sprint(spaces))
			}
			if r >= 0 {
				placement = append(placement, "/")
			}
			spaces = 0
		}
	}

	return assembleFEN(board, strings.Join(placement, "")), true
}
func (board *Board) String() string {
	hdiv := "+---+---+---+---+---+---+---+---+\n"
	buf := bytes.Buffer{}
	buf.WriteString(hdiv)

	for r := 7; r >= 0; r-- {
		rank := r
		if board.Orientation == Black {
			rank = 7 - r
		}

		for f := 0; f < 8; f++ {
			chars := []byte("|   ")
			file := f
			if board.Orientation == Black {
				file = 7 - f
			}

			p := board.squares[rank*8+file]
			if pstr := p.String(); pstr != "" {
				chars[2] = pstr[0]
			}

			buf.Write(chars)
		}

		buf.WriteString(fmt.Sprintf("|  %d\n", r+1))
		buf.WriteString(hdiv)
	}

	if board.Orientation == Black {
		buf.WriteString("  h   g   f   e   d   c   b   a\n")
	} else {
		buf.WriteString("  a   b   c   d   e   f   g   h\n")
	}
	return buf.String()
}

/*
 *
 *
 * FEN structure methods / functions
 */
func (fen FEN) ToBoard() (*Board, error) {
	board := EmptyBoard()
	parts := strings.Split(string(fen), " ")
	if len(parts) != 6 {
		return nil, fmt.Errorf("not enough parts: %d, want 6", len(parts))
	}

	f, r := 0, 7
	for _, symbol := range parts[0] {
		if unicode.IsDigit(symbol) {
			f += int(symbol - '0')
		} else {
			piece := NewPiece(symbol)
			if piece != NoPiece {
				board.squares[r*8+f] = piece
				f++
			} else if symbol != '/' {
				return nil, fmt.Errorf("unknown piece placement symbol: %q", symbol)
			}
		}

		if f >= 8 {
			f = 0
			r--
		}
	}
	if r >= 0 {
		return nil, fmt.Errorf("not enough piece placement symbols")
	}

	switch parts[1] {
	case "w":
		board.SideToMove = White
	case "b":
		board.SideToMove = Black
	default:
		return nil, fmt.Errorf("unknown side to move: %q", parts[1])
	}

	if len(parts[2]) > 4 || len(parts[2]) < 1 {
		return nil, fmt.Errorf("invalid castle character count: %d, want 1-4", len(parts[2]))
	}
	if parts[2] != "-" {
		for _, symbol := range parts[2] {
			switch {
			case symbol == 'K' && !board.CanCastle[WKingside]:
				board.CanCastle[WKingside] = true
			case symbol == 'Q' && !board.CanCastle[WQueenside]:
				board.CanCastle[WQueenside] = true
			case symbol == 'k' && !board.CanCastle[BKingside]:
				board.CanCastle[BKingside] = true
			case symbol == 'q' && !board.CanCastle[BQueenside]:
				board.CanCastle[BQueenside] = true
			default:
				return nil, fmt.Errorf("invalid castle character: %q, want \"|Q|k|q\"", symbol)
			}
		}
	}

	if parts[3] != "-" {
		if len(parts[3]) != 2 {
			return nil, fmt.Errorf("invalid en passant character count: %d, want 2 when not \"-\"", len(parts[3]))
		}

		coord := NewCoord(parts[3])
		if !coord.IsValid() || (coord.Rank() != 2 && coord.Rank() != 5) {
			return nil, fmt.Errorf("invalid en passant coord: [%01b] [%03b] [%03b] [%01b], want [0] [0..111] [0..111] [1]", (uint8(coord)&0x80)>>7, coord.File(), coord.Rank(), coord&1)
		}

		board.EnPassantTarget = coord
	}

	if num, err := strconv.Atoi(parts[4]); err != nil || num < 0 || num > 50 {
		return nil, fmt.Errorf("invalid halfmove clock: %v, want [0..50)", parts[4])
	} else {
		board.HalfmoveClock = num
	}

	if num, err := strconv.Atoi(parts[5]); err != nil || num < 1 {
		return nil, fmt.Errorf("invalid fullmove counter: %v, want [1..]", parts[5])
	} else {
		board.FullmoveCounter = num
	}

	return board, nil
}
