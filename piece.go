package chess

type PieceName int

const (
	Pawn PieceName = iota + 1
	Knight
	Bishop
	Rook
	Queen
	King
)

func NewPieceName(p string) PieceName {
	return map[string]PieceName{
		"N": Knight,
		"B": Bishop,
		"R": Rook,
		"Q": Queen,
		"K": King,
	}[p]
}

func (n PieceName) IsValid() bool {
	return n > 0 && n <= King
}

func (name PieceName) String() string {
	switch name {
	case 0:
		return "none"
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
		return ""
	}
}

type Piece struct {
	Color SideColor
	Name  PieceName
}

func NewPiece(abbr string) Piece {
	return map[string]Piece{
		"p": {Black, Pawn},
		"n": {Black, Knight},
		"b": {Black, Bishop},
		"r": {Black, Rook},
		"q": {Black, Queen},
		"k": {Black, King},

		"P": {White, Pawn},
		"N": {White, Knight},
		"B": {White, Bishop},
		"R": {White, Rook},
		"Q": {White, Queen},
		"K": {White, King},
	}[abbr]
}

func (p Piece) IsValid() bool {
	return p.Color.IsValid() && p.Name.IsValid()
}

func (piece Piece) String() string {
	if !piece.IsValid() {
		return ""
	}

	name := [...]rune{'p', 'n', 'b', 'r', 'q', 'k'}[int(piece.Name)-1]
	if piece.Color == White {
		name += 'A' - 'a' // capitalize name
	}
	return string(name)
}
