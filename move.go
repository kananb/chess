package chess

import (
	"bytes"
	"fmt"
	"regexp"
)

type Coord struct {
	File, Rank int
}

func NewCoord(c string) Coord {
	if len(c) != 2 || c[0] < 'a' || c[0] > 'h' || c[1] < '1' || c[1] > '8' {
		return Coord{0, 0}
	}

	return Coord{int(c[0]-'a') + 1, int(c[1]-'1') + 1}
}
func indexCoord(i int) Coord {
	return Coord{i&7 + 1, i/8 + 1}
}

func (c Coord) Index() int {
	return c.Rank*8 + c.File - 9
}
func (c Coord) IsValid() bool {
	return c.File > 0 && c.File <= 8 && c.Rank > 0 && c.Rank <= 8
}

func (c Coord) String() string {
	if c.File == 0 || c.Rank == 0 {
		return ""
	}
	return fmt.Sprintf("%c%d", 'a'+rune(c.File-1), c.Rank)
}

type CheckType int

const (
	Check CheckType = iota + 1
	Checkmate
)

func (c CheckType) IsValid() bool {
	return c == Check || c == Checkmate
}

type MoveFlags struct {
	Moves       PieceName
	Captures    PieceName
	PromotesTo  PieceName
	CastlesTo   CastleSide
	Check       CheckType
	IsEnPassant bool
	OffersDraw  bool
}

// Move structure
type Move struct {
	From, To Coord
	MoveFlags
}

var moveexp = regexp.MustCompile(`^(?P<Move>(?P<Piece>[NBRQK]?)(?:(?P<FileSpecifier>[a-h])?(?P<RankSpecifier>[1-8])?)(?P<Takes>x?)(?P<Destination>[a-h][1-8])(?:=?(?P<Promotion>[NBRQ]))?|(?P<Castle>[0O](?:-[0O]){1,2}))(?P<Check>\+{0,2}|#)(?P<EnPassant> e\.p\.)?(?P<OffersDraw> \(=\))?$`)

func NewMove(san string, board *Board) (move Move, err error) {
	matches := moveexp.FindStringSubmatch(san)
	if matches == nil {
		return move, fmt.Errorf("invalid move")
	}

	// Draw offer
	if i := moveexp.SubexpIndex("OffersDraw"); i != -1 && matches[i] != "" {
		move.OffersDraw = true
	}

	// En passant
	if i := moveexp.SubexpIndex("EnPassant"); i != -1 && matches[i] != "" {
		move.IsEnPassant = true
	}

	// Check
	if i := moveexp.SubexpIndex("Check"); i != -1 && matches[i] != "" {
		if matches[i] == "+" || matches[i] == "++" {
			move.Check = Check
		} else if matches[i] == "#" {
			move.Check = Checkmate
		} else {
			panic(fmt.Sprintf("invalid check notation %q, want +|++|#", matches[i]))
		}
	}

	// Castling
	if i := moveexp.SubexpIndex("Castle"); i != -1 && matches[i] != "" {
		rank := 1
		if board.SideToMove == Black {
			rank = 8
		}
		move.From = Coord{5, rank}

		if matches[i] == "O-O" || matches[i] == "0-0" {
			move.CastlesTo = Kingside
			move.To = Coord{7, rank}
		} else if matches[i] == "O-O-O" || matches[i] == "0-0-0" {
			move.CastlesTo = Queenside
			move.To = Coord{3, rank}
		} else {
			panic(fmt.Sprintf("invalid castle notation %q", matches[i]))
		}

		move.Moves = King
		return // no need to continue parsing
	}

	// Pawn promotion
	if i := moveexp.SubexpIndex("Promotion"); i != -1 && matches[i] != "" {
		pieceName := NewPieceName(string(matches[i][0]))
		if pieceName == 0 || pieceName >= King {
			panic(fmt.Sprintf("invalid promotion notation %q, want [NBRQ]", matches[i]))
		}

		move.PromotesTo = pieceName
	}

	// Destination
	if i := moveexp.SubexpIndex("Destination"); i != -1 && matches[i] != "" {
		if dest := NewCoord(matches[i]); dest.File == 0 || dest.Rank == 0 {
			panic(fmt.Sprintf("invalid destination notation %q, want [a-h][1-8]", matches[i]))
		} else {
			move.To = dest
		}
	} else {
		return move, fmt.Errorf("Move must include destination square")
	}

	// Takes
	if i := moveexp.SubexpIndex("Takes"); i != -1 && matches[i] != "" {
		pieceName := board.At(move.To).Name
		move.Captures = pieceName
	}

	// Move piece
	if i := moveexp.SubexpIndex("Piece"); i != -1 && matches[i] != "" {
		move.Moves = NewPieceName(string(matches[i][0]))
		if !move.Moves.IsValid() {
			panic(fmt.Sprintf("invalid move piece notation %q, want [NBRQK]?", matches[i]))
		}
	} else {
		move.Moves = Pawn
	}

	// Piece disambiguation
	file, rank := -1, -1
	if i := moveexp.SubexpIndex("FileSpecifier"); i != -1 && matches[i] != "" {
		file = int(matches[i][0] - 'a' + 1)
	}
	if i := moveexp.SubexpIndex("RankSpecifier"); i != -1 && matches[i] != "" {
		rank = int(matches[i][0] - '1' + 1)
	}
	move.From = Coord{file, rank}

	if file == -1 || rank == -1 {
		candidates := board.PseudoMoves(move.Moves)

		for _, c := range candidates {
			if c.To == move.To && (file == -1 || c.From.File == file) && (rank == -1 || c.From.Rank == rank) {
				if move.From.IsValid() {
					return move, fmt.Errorf("Move is ambiguous")
				}
				move.From = c.From
			}
		}
	}

	if !move.From.IsValid() {
		return move, fmt.Errorf("invalid move")
	} else if board.At(move.From).Name != move.Moves {
		return move, fmt.Errorf("wrong piece type")
	}

	return
}

func (m Move) Matches(move Move) bool {
	return m.To == move.To && m.From == move.From
}
func (m Move) IsValid() bool {
	return m.To.IsValid() && m.From.IsValid()
}

func (m Move) String() string {
	if !m.IsValid() {
		return ""
	}

	buf := bytes.Buffer{}

	if m.CastlesTo == Kingside {
		buf.WriteString("O-O")
	} else if m.CastlesTo == Queenside {
		buf.WriteString("O-O-O")
	} else {
		buf.WriteString(m.Moves.Abbreviation())
		buf.WriteString(m.From.String())
		if m.Captures.IsValid() {
			buf.WriteString("x")
		}
		buf.WriteString(m.To.String())
		if m.PromotesTo.IsValid() {
			buf.WriteString("=" + m.PromotesTo.Abbreviation())
		}
	}

	if m.Check == Checkmate {
		buf.WriteString("#")
	} else if m.Check == Check {
		buf.WriteString("+")
	}
	if !m.CastlesTo.IsValid() && m.IsEnPassant {
		buf.WriteString(" e.p.")
	}
	if m.OffersDraw {
		buf.WriteString(" (=)")
	}

	return buf.String()
}
