package chess

import (
	"fmt"
)

var directionOffets = [8]int8{8, -8, 1, -1, 9, -7, -9, 7}
var squaresToEdge = [64][8]uint8{}

const (
	north int8 = iota
	south
	east
	west
	northeast
	southeast
	southwest
	northwest
)

func init() {
	min := func(a, b uint8) uint8 {
		if a < b {
			return a
		} else {
			return b
		}
	}

	for r := uint8(0); r < 8; r++ {
		for f := uint8(0); f < 8; f++ {
			pos := squaresToEdge[r*8+f]
			pos[north] = 7 - r
			pos[south] = r
			pos[east] = 7 - f
			pos[west] = f
			pos[northeast] = min(pos[north], pos[east])
			pos[southeast] = min(pos[south], pos[east])
			pos[southwest] = min(pos[south], pos[west])
			pos[northwest] = min(pos[north], pos[west])
		}
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
func capturesFriendly(start, target Piece) bool {
	return target != NoPiece && target.Color() == start.Color()
}
func (move Move) isPawnMove(board *Board) bool {
	start, target := board.Get(move.Start), board.Get(move.Target)
	canMoveDouble := (start.Color() == White && move.Start.Rank() == 1) || (start.Color() == Black && move.Start.Rank() == 6)
	dx, dy := move.Translation()
	adx, ady := abs(dx), abs(dy)

	if adx > 1 || ady < 1 || ady > 2 || (ady == 2 && (adx > 0 || !canMoveDouble)) {
		println("impossible movement", adx, ady)
		return false
	}

	if (dy < 0 && start.Color() != Black) || (dy > 0 && start.Color() != White) {
		println("backwards movement")
		return false
	}

	if dx != 0 {
		if target == NoPiece && board.EnPassantTarget != move.Target {
			println("capture on empty square")
			return false
		} else if capturesFriendly(start, target) {
			println("capture friendly piece")
			return false
		}
	} else if target != NoPiece || (ady == 2 && board.squares[move.Start.Rank()*8+move.Start.File()+4*dy] != NoPiece) {
		println("blocked ahead")
		return false
	}

	return true
}
func (move Move) isKnightMove(board *Board) bool {
	start, target := board.Get(move.Start), board.Get(move.Target)
	dx, dy := move.Translation()
	adx, ady := abs(dx), abs(dy)

	return adx != 0 && ady != 0 && adx+ady == 3 && !capturesFriendly(start, target)
}
func (move Move) isSlidingMove(board *Board) bool {
	return true
}
func (move Move) isKingMove(board *Board) bool {
	start, target := board.Get(move.Start), board.Get(move.Target)
	dx, dy := move.Translation()
	adx, ady := abs(dx), abs(dy)

	return adx < 2 && ady < 2 && adx+ady > 0 && !capturesFriendly(start, target)
}
func (move Move) isPseudoLegal(board *Board) bool {
	if !move.IsValid(board) || board.Get(move.Start).Color() != board.SideToMove {
		return false
	}

	switch board.Get(move.Start).Type() {
	case Pawn:
		return move.isPawnMove(board)
	case Knight:
		return move.isKnightMove(board)
	case Bishop, Rook, Queen:
		return move.isSlidingMove(board)
	case King:
		return move.isKingMove(board)
	default:
		panic(fmt.Sprintf("unknown piece type: %d", board.Get(move.Start).Type()))
	}
}
func (move Move) isLegal(board *Board) bool {
	if !move.isPseudoLegal(board) {
		return false
	}

	return true
}
func (board *Board) MakeMove(move Move) bool {
	if !move.isLegal(board) {
		return false
	}

	piece := board.Get(move.Start)
	capture := board.Get(move.Target) != NoPiece
	pawnMove := piece.Type() == Pawn

	board.Set(move.Target, piece)
	board.Set(move.Start, NoPiece)

	board.SideToMove ^= 1
	if piece.Color() == Black {
		board.FullmoveCounter++
	}
	if !capture && !pawnMove {
		board.HalfmoveClock++
	} else {
		board.HalfmoveClock = 0
	}

	return true
}
func (board *Board) UnmakeMove(move Move) bool {
	// TODO
	return false
}
