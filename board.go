package chess

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type PieceType uint8

const ( // piece types
	NoType PieceType = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King

	typeMask  uint8 = 0b111
	typeShift uint8 = 2
)

func (pt PieceType) String() string {
	switch pt {
	case Pawn:
		return "pawn"
	case Knight:
		return "knight"
	case Bishop:
		return "bishop"
	case Rook:
		return "rook"
	case Queen:
		return "queen"
	case King:
		return "king"
	default:
		return "none"
	}
}

//
// --------------------------------------------------------
// PIECE TYPE
//
//
//
// COLOR
// --------------------------------------------------------
//

type SideColor uint8

const ( // side colors
	NoColor SideColor = iota
	White
	Black

	colorMask  uint8 = 0b11
	colorShift uint8 = 0
)

func (sc SideColor) String() string {
	if sc == White {
		return "w"
	} else if sc == Black {
		return "b"
	}
	return "-"
}

//
// --------------------------------------------------------
// COLOR
//
//
//
// PIECE
// --------------------------------------------------------
//

// Piece encoding
// msb  [0000]  [000..110] [0..1] lsb
//      unused     type     color
//
// Example (black rook):	0000 100 1
// Example (white bishop):	0000 011 0
type Piece uint8

const NoPiece = Piece(0)

func NewPiece(c SideColor, t PieceType) Piece {
	if t > King {
		return Piece(0)
	}
	return Piece((uint8(t)&typeMask)<<typeShift | (uint8(c)&colorMask)<<colorShift)
}
func PieceFromRune(r rune) (Piece, bool) {
	p, ok := map[rune]Piece{
		'p': NewPiece(Black, Pawn),
		'n': NewPiece(Black, Knight),
		'b': NewPiece(Black, Bishop),
		'r': NewPiece(Black, Rook),
		'q': NewPiece(Black, Queen),
		'k': NewPiece(Black, King),

		'P': NewPiece(White, Pawn),
		'N': NewPiece(White, Knight),
		'B': NewPiece(White, Bishop),
		'R': NewPiece(White, Rook),
		'Q': NewPiece(White, Queen),
		'K': NewPiece(White, King),
	}[r]

	return p, ok
}

func (p Piece) Type() PieceType {
	return PieceType(uint8(p) >> typeShift & typeMask)
}
func (p Piece) Color() SideColor {
	return SideColor(uint8(p) >> colorShift & colorMask)
}
func (piece Piece) String() string {
	if piece.Type() > King || piece == NoPiece {
		return ""
	}

	name := [...]rune{'p', 'n', 'b', 'r', 'q', 'k'}[int(piece.Type())-1]
	if piece.Color() == White {
		name += 'A' - 'a' // capitalize name
	}
	return string(name)
}

//
// --------------------------------------------------------
// PIECE
//
//
//
// COORD
// --------------------------------------------------------
//

// Square encoding:
// msb  [00]   [000..111] [000..111]  lsb
//    is_valid    file       rank
//
// Example (b6):	11 001 101
// Example (h2):	11 111 001
// Example (a1):	11 000 000

// The zero value represents an empty coordinate, but if its first bit is a 1, then it is the coordinate (0, 0) (a1)
type Coord uint8

const (
	NoCoord Coord = Coord(0)

	coordMask uint8 = 0b111
	indexMask uint8 = 0b111111
)

func NewCoord(f, r int) Coord {
	if f < 0 || f > 7 || r < 0 || r > 7 {
		return NoCoord
	}
	return Coord(1<<7 | uint8(r)<<3 | uint8(f))
}
func CoordFromString(s string) (Coord, bool) {
	if len(s) != 2 || s[0] < 'a' || s[0] > 'h' || s[1] < '1' || s[1] > '8' {
		return Coord(0), false
	}

	return NewCoord(int(s[0]-'a'), int(s[1]-'1')), true
}
func coordFromIndex(i int) Coord {
	if i < 0 || i > 63 {
		return NoCoord
	}
	return Coord(1<<7 | uint8(i))
}

func (c Coord) File() int {
	return int(uint8(c) & coordMask) // mod by 8
}
func (c Coord) Rank() int {
	return int(uint8(c) >> 3 & coordMask) // divide by 8
}
func (c Coord) index() int {
	return int(uint8(c) & indexMask)
}

func (c Coord) String() string {
	return fmt.Sprintf("%c%d", 'a'+rune(c.File()), c.Rank()+1)
}

//
// --------------------------------------------------------
// COORD
//
//
//
// MOVE
// --------------------------------------------------------
//

// Move structure
type Move struct {
	Start, End Coord
}

func NewMove(start, end Coord) Move {
	return Move{start, end}
}

func (m Move) String() string {
	return fmt.Sprintf("%v->%v", m.Start, m.End)
}

//
// --------------------------------------------------------
// MOVE
//
//
//
// CASTLES
// --------------------------------------------------------
//

type CastleSide uint8

const (
	Kingside CastleSide = iota
	Queenside
)

// Castle Structure
type Castles [4]bool

func NewCastles() (castles Castles) {
	return
}
func CastlesFromString(s string) (castles Castles, ok bool) {
	for _, symbol := range s {
		switch {
		case symbol == 'K' && !castles.Can(White, Kingside):
			castles.Allow(White, Kingside)
		case symbol == 'Q' && !castles.Can(White, Queenside):
			castles.Allow(White, Queenside)
		case symbol == 'k' && !castles.Can(Black, Kingside):
			castles.Allow(Black, Kingside)
		case symbol == 'q' && !castles.Can(Black, Queenside):
			castles.Allow(Black, Queenside)
		default:
			return castles, false
		}
	}

	return castles, true
}

func (c *Castles) Allow(sc SideColor, side CastleSide) {
	if (sc != White && sc != Black) || (side != Kingside && side != Queenside) {
		return
	}

	c[uint(sc&0b10)|uint(side)] = true
}
func (c *Castles) Can(sc SideColor, side CastleSide) bool {
	if (sc != White && sc != Black) || (side != Kingside && side != Queenside) {
		return false
	}

	return c[uint(sc&0b10)|uint(side)]
}

func (c *Castles) String() string {
	buf := bytes.Buffer{}

	if c.Can(White, Kingside) {
		buf.WriteByte('K')
	}
	if c.Can(White, Queenside) {
		buf.WriteByte('Q')
	}
	if c.Can(Black, Kingside) {
		buf.WriteByte('k')
	}
	if c.Can(Black, Queenside) {
		buf.WriteByte('q')
	}

	if buf.Len() == 0 {
		return "-"
	}
	return buf.String()
}

//
// --------------------------------------------------------
// CASTLES
//
//
//
// BOARD
// --------------------------------------------------------
//

// A chess board structure
// Keeps track of piece positions, board orientation,
// side-to-move, castling ability, en passant target,
// the halfmove clock, and fullmove counter
type Board struct {
	squares     [64]Piece
	Orientation SideColor

	SideToMove      SideColor
	CastleRights    Castles
	EnPassantTarget Coord
	HalfmoveClock   int
	FullmoveCounter int
}

func NewBoard() *Board {
	board := new(Board)
	board.Orientation = White
	board.SideToMove = White
	board.FullmoveCounter = 1

	return board
}
func BoardFromString(fen string) (*Board, error) {
	fields := strings.Split(string(fen), " ")
	if len(fields) != 6 {
		return nil, fmt.Errorf("not enough fields: %d, want 6", len(fields))
	}

	board := NewBoard()

	// Extract piece placement information
	f, r := 0, 7
	for _, symbol := range fields[0] {
		if unicode.IsDigit(symbol) {
			f += int(symbol - '0')
		} else {
			piece, ok := PieceFromRune(symbol)
			if ok {
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

	// Extract side to move
	switch fields[1] {
	case "w":
		board.SideToMove = White
	case "b":
		board.SideToMove = Black
	default:
		return nil, fmt.Errorf("unknown side to move: %q", fields[1])
	}

	// Extract castling ability
	if len(fields[2]) > 4 || len(fields[2]) < 1 {
		return nil, fmt.Errorf("invalid castle character count: %d, want 1-4", len(fields[2]))
	}
	if fields[2] != "-" {
		if castles, ok := CastlesFromString(fields[2]); !ok {
			return nil, fmt.Errorf("invalid castle string: %q, want \"-|K?Q?k?q?\"", fields[2])
		} else {
			board.CastleRights = castles
		}
	}

	// Extract en passant square
	if fields[3] != "-" {
		if coord, ok := CoordFromString(fields[3]); !ok {
			return nil, fmt.Errorf("invalid en passant string: %q, want \"([a-h])([1-8])\"", fields[3])
		} else {
			board.EnPassantTarget = coord
		}
	}

	// Extract halfmove clock
	if num, err := strconv.Atoi(fields[4]); err != nil || num < 0 || num > 50 {
		return nil, fmt.Errorf("invalid halfmove clock: %v, want [0..50)", fields[4])
	} else {
		board.HalfmoveClock = num
	}

	// Extract fullmove counter
	if num, err := strconv.Atoi(fields[5]); err != nil || num < 1 {
		return nil, fmt.Errorf("invalid fullmove counter: %v, want [1..]", fields[5])
	} else {
		board.FullmoveCounter = num
	}

	return board, nil
}

func (board *Board) At(c Coord) *Piece {
	if c == NoCoord {
		return nil
	}
	return &board.squares[c.index()]
}
func (board *Board) Flip() {
	if board.Orientation == White {
		board.Orientation = Black
	} else {
		board.Orientation = White
	}
}
func (board *Board) FEN() string {
	placement := bytes.Buffer{}

	spaces := 0
	for f, r := 0, 7; r >= 0; {
		piece := board.squares[r*8+f]
		if piece == NoPiece {
			spaces++
		} else {
			if pieceString := piece.String(); pieceString == "" {
				panic(fmt.Sprintf("unknown piece value found: %v", piece))
			} else {
				if spaces > 0 {
					placement.WriteString(fmt.Sprint(spaces))
				}
				placement.WriteString(pieceString)
				spaces = 0
			}
		}

		f++
		if f&7 == 0 { // file is divisible by 8
			r -= 1
			f = 0

			if spaces > 0 {
				placement.WriteString(fmt.Sprint(spaces))
			}
			if r >= 0 {
				placement.WriteByte('/')
			}
			spaces = 0
		}
	}

	epTarget := "-"
	if board.EnPassantTarget != Coord(0) {
		epTarget = board.EnPassantTarget.String()
	}

	return fmt.Sprintf("%v %v %v %v %v %v", placement.String(), board.SideToMove, board.CastleRights.String(), epTarget, board.HalfmoveClock, board.FullmoveCounter)
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

		buf.WriteString(fmt.Sprintf("|  %d\n", rank+1))
		buf.WriteString(hdiv)
	}

	if board.Orientation == White {
		buf.WriteString("  a   b   c   d   e   f   g   h\n")
	} else {
		buf.WriteString("  h   g   f   e   d   c   b   a\n")
	}
	return buf.String()
}

func StartingPosition() *Board {
	return &Board{
		squares: [64]Piece{
			NewPiece(White, Rook), NewPiece(White, Knight), NewPiece(White, Bishop), NewPiece(White, Queen), NewPiece(White, King), NewPiece(White, Bishop), NewPiece(White, Knight), NewPiece(White, Rook),
			NewPiece(White, Pawn), NewPiece(White, Pawn), NewPiece(White, Pawn), NewPiece(White, Pawn), NewPiece(White, Pawn), NewPiece(White, Pawn), NewPiece(White, Pawn), NewPiece(White, Pawn),
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			NewPiece(Black, Pawn), NewPiece(Black, Pawn), NewPiece(Black, Pawn), NewPiece(Black, Pawn), NewPiece(Black, Pawn), NewPiece(Black, Pawn), NewPiece(Black, Pawn), NewPiece(Black, Pawn),
			NewPiece(Black, Rook), NewPiece(Black, Knight), NewPiece(Black, Bishop), NewPiece(Black, Queen), NewPiece(Black, King), NewPiece(Black, Bishop), NewPiece(Black, Knight), NewPiece(Black, Rook),
		},
		Orientation:     White,
		SideToMove:      White,
		CastleRights:    [4]bool{true, true, true, true},
		FullmoveCounter: 1,
	}
}
