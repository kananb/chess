package chess

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"unicode"
)

type SideColor int

const ( // side colors
	White SideColor = iota + 1
	Black
)

func NewSideColor(c string) SideColor {
	return map[string]SideColor{
		"w": White,
		"b": Black,
	}[c]
}

func (c SideColor) IsValid() bool {
	return c == White || c == Black
}

func (c SideColor) String() string {
	if c == White {
		return "w"
	} else if c == Black {
		return "b"
	}
	return ""
}

type CastleSide int

const (
	Kingside CastleSide = iota + 1
	Queenside
)

func (s CastleSide) IsValid() bool {
	return s == Kingside || s == Queenside
}

// Castle Structure
type Castles int

func NewCastles(s string) (c Castles) {
	for _, symbol := range s {
		switch {
		case symbol == 'K' && !c.Can(White, Kingside):
			c.Allow(White, Kingside)
		case symbol == 'Q' && !c.Can(White, Queenside):
			c.Allow(White, Queenside)
		case symbol == 'k' && !c.Can(Black, Kingside):
			c.Allow(Black, Kingside)
		case symbol == 'q' && !c.Can(Black, Queenside):
			c.Allow(Black, Queenside)
		default:
			return 0
		}
	}

	return c
}

func (c *Castles) Allow(color SideColor, side CastleSide) {
	if !color.IsValid() || !side.IsValid() {
		return
	}

	shift := (uint8(color)-1)*2 + uint8(side) - 1
	*c |= 1 << shift
}
func (c *Castles) Disallow(color SideColor, side CastleSide) {
	if !color.IsValid() || !side.IsValid() {
		return
	}

	mask := Castles(^(uint8(side) << ((uint8(color) - 1) * 2)) & 15)
	*c = *c & mask
}
func (c *Castles) Can(color SideColor, side CastleSide) bool {
	if !color.IsValid() || !side.IsValid() {
		return false
	}

	mask := Castles((uint8(side) << ((uint8(color) - 1) * 2)) & 15)
	return *c&mask != 0
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

type BoardData struct {
	SideToMove      SideColor
	CastleRights    Castles
	EnPassantTarget Coord
	HalfmoveClock   int
	FullmoveCounter int
}

type BoardState struct {
	Move
	BoardData
}

// A chess board structure
// Keeps track of piece positions, board orientation,
// side-to-move, castling ability, en passant target,
// the halfmove clock, and fullmove counter
type Board struct {
	squares [64]Piece
	BoardData

	history []BoardState
}

var fenexp = regexp.MustCompile(`^(?P<PiecePlacement>(?:[pnbrqkPNBRQK1-8]{1,8}\/){7}[pnbrqkPNBRQK1-8]{1,8})\s+(?P<SideToMove>b|w)\s+(?P<Castling>-|K?Q?k?q?)\s+(?P<EnPassant>-|[a-h][3-6])\s+(?P<HalfmoveClock>\d+)\s+(?P<FullmoveCounter>\d+)\s*$`)

func NewBoard(fen string) (board *Board, err error) {
	board = new(Board)
	board.SideToMove = White
	board.FullmoveCounter = 1
	board.history = make([]BoardState, 0, 128)

	if fen == "" {
		return
	}

	matches := fenexp.FindStringSubmatch(fen)
	if matches == nil {
		return nil, fmt.Errorf("invalid FEN string")
	}

	if i := fenexp.SubexpIndex("PiecePlacement"); i != -1 && matches[i] != "" {
		f, r := 0, 7
		for _, symbol := range matches[i] {
			if unicode.IsDigit(symbol) {
				f += int(symbol - '0')
			} else {
				piece := NewPiece(string(symbol))
				if piece.IsValid() {
					board.squares[r*8+f] = piece
					f++
				} else if symbol != '/' {
					return nil, fmt.Errorf("unknown piece symbol: %q", symbol)
				}
			}

			if f >= 8 {
				f = 0
				r--
			}
		}
		if r >= 0 {
			return nil, fmt.Errorf("not enough piece symbols")
		}
	} else {
		return nil, fmt.Errorf("invalid or missing piece placement")
	}

	if i := fenexp.SubexpIndex("SideToMove"); i != -1 && matches[i] != "" {
		if matches[i] == "b" {
			board.SideToMove = Black
		}
	} else {
		return nil, fmt.Errorf("invalid or missing side to move")
	}

	if i := fenexp.SubexpIndex("Castling"); i != -1 && matches[i] != "" {
		board.CastleRights = NewCastles(matches[i])
	} else {
		return nil, fmt.Errorf("invalid or missing castling rights")
	}

	if i := fenexp.SubexpIndex("EnPassant"); i != -1 && matches[i] != "" {
		board.EnPassantTarget = NewCoord(matches[i])
	} else {
		return nil, fmt.Errorf("invalid or missing en passant")
	}

	if i := fenexp.SubexpIndex("HalfmoveClock"); i != -1 && matches[i] != "" {
		if ply, err := strconv.Atoi(matches[i]); err == nil {
			if ply < 0 || ply > 50 {
				return nil, fmt.Errorf("halfmove clock out of range, [0, 50]")
			}
			board.HalfmoveClock = ply
		} else {
			panic("invalid halfmove clock matched regex")
		}
	}

	if i := fenexp.SubexpIndex("FullmoveCounter"); i != -1 && matches[i] != "" {
		if counter, err := strconv.Atoi(matches[i]); err == nil {
			if counter < 0 {
				return nil, fmt.Errorf("fullmove counter out of range, [0, inf]")
			}
			board.FullmoveCounter = counter
		} else {
			panic("invalid fullmove counter matched regex")
		}
	} else {
		board.FullmoveCounter = 1
	}

	return
}

func (board *Board) At(c Coord) *Piece {
	if !c.IsValid() {
		return nil
	}
	return &board.squares[c.Index()]
}
func (board *Board) History() []string {
	history := make([]string, 0, len(board.history))

	for _, state := range board.history {
		history = append(history, state.Move.String())
	}

	return history
}
func (board *Board) pieceIndices(side SideColor, names ...PieceName) []int {
	pieces := make([]int, 0, 16)

	for i := 0; i < len(board.squares); i++ {
		if board.squares[i].Color != side {
			continue
		}
		if len(names) == 0 {
			pieces = append(pieces, i)
			continue
		}

		for _, name := range names {
			if board.squares[i].Name == name {
				pieces = append(pieces, i)
				break
			}
		}
	}

	return pieces
}

func (board *Board) Ascii() string {
	hdiv := "+---+---+---+---+---+---+---+---+\n"
	buf := bytes.Buffer{}
	buf.WriteString(hdiv)

	for r := 7; r >= 0; r-- {

		for f := 0; f < 8; f++ {
			chars := []byte("|   ")

			p := board.squares[r*8+f]
			if pstr := p.String(); pstr != "" {
				chars[2] = pstr[0]
			}

			buf.Write(chars)
		}

		buf.WriteString(fmt.Sprintf("|  %d\n", r+1))
		buf.WriteString(hdiv)
	}

	buf.WriteString("  a   b   c   d   e   f   g   h\n")
	return buf.String()
}
func (board *Board) String() string {
	placement := bytes.Buffer{}

	spaces := 0
	for f, r := 0, 7; r >= 0; {
		piece := board.squares[r*8+f]
		if piece.Name == 0 {
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
	if board.EnPassantTarget.IsValid() {
		epTarget = board.EnPassantTarget.String()
	}

	return fmt.Sprintf("%v %v %v %v %v %v", placement.String(), board.SideToMove, board.CastleRights.String(), epTarget, board.HalfmoveClock, board.FullmoveCounter)
}

func StartingPosition() *Board {
	return &Board{
		squares: [64]Piece{
			{White, Rook}, {White, Knight}, {White, Bishop}, {White, Queen}, {White, King}, {White, Bishop}, {White, Knight}, {White, Rook},
			{White, Pawn}, {White, Pawn}, {White, Pawn}, {White, Pawn}, {White, Pawn}, {White, Pawn}, {White, Pawn}, {White, Pawn},
			{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
			{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
			{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
			{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0},
			{Black, Pawn}, {Black, Pawn}, {Black, Pawn}, {Black, Pawn}, {Black, Pawn}, {Black, Pawn}, {Black, Pawn}, {Black, Pawn},
			{Black, Rook}, {Black, Knight}, {Black, Bishop}, {Black, Queen}, {Black, King}, {Black, Bishop}, {Black, Knight}, {Black, Rook},
		},
		BoardData: BoardData{
			SideToMove:      White,
			CastleRights:    15,
			FullmoveCounter: 1,
		},
	}
}
