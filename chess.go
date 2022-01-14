package chess

/* Piece encoding
 * msb [0000]  [000..110] [0..1] lsb
 *     unused     type     color
 *
 * Example (black rook):	0000 100 1
 * Example (white bishop):	0000 011 0
 */
type Piece uint8

/* Square encoding
 * msb  [0]   [000..111] [000..111]  [0..1]  lsb
 *     unused    file       rank    no-coord
 *
 * Example (b6):	0 001 110 1
 * Example (h2):	0 111 010 1
 * Example (none):	0 000 000 0
 */
type Coord uint8

type Move struct {
	Start, Target Coord
}

const ( // piece types
	NoPiece Piece = iota

	Pawn uint8 = iota << 1
	Knight
	Bishop
	Rook
	Queen
	King

	typeMask  uint8 = 0b1110
	typeShift uint8 = 1
)
const ( // side colors
	White uint8 = iota
	Black

	colorMask  uint8 = 0b1
	colorShift uint8 = 0
)
const (
	NoCoord Coord = 0

	File      uint8 = 0b1110000
	fileShift uint8 = 4
	Rank      uint8 = 0b1110
	rankShift uint8 = 1
)

type Board struct {
	squares         [64]Piece
	SideToMove      uint8
	CanCastle       [4]bool
	EnPassantTarget Coord
	HalfmoveClock   int
	FullmoveCounter int
	Orientation     uint8
}
type FEN string

const (
	WKingside int = iota
	WQueenside
	BKingside
	BQueenside
)
