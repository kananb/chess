package chess

import "fmt"

var slideDirections = [...]struct{ f, r int }{
	{1, 1}, {1, -1}, {-1, -1}, {-1, 1},
	{0, 1}, {1, 0}, {0, -1}, {-1, 0},
}
var knightOffsets = [...]struct{ f, r int }{
	{1, 2}, {2, 1}, {2, -1}, {1, -2},
	{-1, -2}, {-2, -1}, {-2, 1}, {-1, 2},
}

var promoteTypes = [...]PieceName{
	Knight, Bishop, Rook, Queen,
}

func (board *Board) pawnMoves(pieces []int) []Move {
	moveSet := make([]Move, 0, 16)

	dir := 1
	if board.SideToMove == Black {
		dir = -1
	}

	for _, i := range pieces {
		piece := board.squares[i]
		if piece.Name != Pawn {
			continue
		}

		from := indexCoord(i)
		to := Coord{from.File, from.Rank + dir}
		if to.IsValid() && !board.At(to).IsValid() {
			if to.Rank == 1 || to.Rank == 8 {
				for _, t := range promoteTypes {
					moveSet = append(moveSet, Move{from, to, MoveFlags{Moves: Pawn, PromotesTo: t}})
				}
			} else {
				moveSet = append(moveSet, Move{from, to, MoveFlags{Moves: Pawn}})
			}

			to = Coord{from.File, from.Rank + dir*2}
			canMoveDouble := (piece.Color == White && from.Rank == 2) || (piece.Color == Black && from.Rank == 7)
			if to.IsValid() && canMoveDouble && !board.At(to).IsValid() {
				moveSet = append(moveSet, Move{from, to, MoveFlags{Moves: Pawn}})
			}
		}

		for off := 3; off > 0; off -= 2 {
			to = Coord{from.File + off - 2, from.Rank + dir}
			if to.IsValid() && ((board.At(to).IsValid() && board.At(to).Color != piece.Color) || board.EnPassantTarget == to) {
				if to.Rank == 1 || to.Rank == 8 {
					for _, t := range promoteTypes {
						moveSet = append(moveSet, Move{from, to, MoveFlags{Moves: Pawn, Captures: board.At(to).Name, PromotesTo: t}})
					}
				} else {
					enPassant := board.EnPassantTarget == to
					moveSet = append(moveSet, Move{from, to, MoveFlags{Moves: Pawn, Captures: board.At(to).Name, IsEnPassant: enPassant}})
				}
			}
		}
	}

	return moveSet
}
func (board *Board) knightMoves(pieces []int) []Move {
	moveSet := make([]Move, 0, 16)

	for _, i := range pieces {
		piece := board.squares[i]
		if piece.Name != Knight {
			continue
		}

		from := indexCoord(i)
		for _, off := range knightOffsets {
			to := Coord{from.File + off.f, from.Rank + off.r}
			if to.IsValid() && board.At(to).Color != piece.Color {
				moveSet = append(moveSet, Move{from, to, MoveFlags{Moves: Knight, Captures: board.At(to).Name}})
			}
		}
	}

	return moveSet
}
func (board *Board) slidingMoves(pieces []int) []Move {
	moveSet := make([]Move, 0, 128)

	for _, i := range pieces {
		piece := board.squares[i]
		if piece.Name != Bishop && piece.Name != Rook && piece.Name != Queen {
			continue
		}

		di, df := 0, 8
		if piece.Name == Bishop {
			df = 4
		} else if piece.Name == Rook {
			di = 4
		}

		from := indexCoord(i)
		for d := di; d < df; d++ {
			for off := 1; ; off++ {
				to := Coord{from.File + slideDirections[d].f*off, from.Rank + slideDirections[d].r*off}
				if !to.IsValid() || board.At(to).Color == piece.Color {
					break
				}

				moveSet = append(moveSet, Move{from, to, MoveFlags{Moves: piece.Name, Captures: board.At(to).Name}})
				if board.At(to).IsValid() {
					break
				}
			}
		}
	}

	return moveSet
}
func (board *Board) kingMoves(pieces []int) []Move {
	moveSet := make([]Move, 0, 8)

	for _, i := range pieces {
		piece := board.squares[i]
		if piece.Name != King {
			continue
		}

		from := indexCoord(i)
		for x := -1; x < 2; x++ {
			for y := -1; y < 2; y++ {
				to := Coord{from.File + x, from.Rank + y}
				if to.IsValid() && board.At(to).Color != piece.Color {
					moveSet = append(moveSet, Move{from, to, MoveFlags{Moves: King, Captures: board.At(to).Name}})
				}
			}
		}

		checkCastle := func(side CastleSide, dir int) {
			between, to := Coord{from.File + 1*dir, from.Rank}, Coord{from.File + 2*dir, from.Rank}
			if board.CastleRights.Can(piece.Color, side) && !board.At(between).IsValid() && !board.At(to).IsValid() {
				moveSet = append(moveSet, Move{from, to, MoveFlags{Moves: King, CastlesTo: side}})
			}
		}
		if (from.Rank == 1 || from.Rank == 8) && from.File == 5 {
			checkCastle(Kingside, 1)
			checkCastle(Queenside, -1)
		}

		break
	}

	return moveSet
}

func (board *Board) InCheck(side SideColor) bool {
	var from Coord
	for i := 0; i < len(board.squares); i++ {
		if board.squares[i].Color == side && board.squares[i].Name == King {
			from = indexCoord(i)
			break
		}
	}

	for i, dir := range slideDirections {
		for off := 1; ; off++ {
			to := Coord{from.File + dir.f*off, from.Rank + dir.r*off}
			if !to.IsValid() || board.At(to).Color == side {
				break
			}

			target := board.At(to).Name
			if target == Queen || (i < 4 && target == Bishop) || (i >= 4 && target == Rook) {
				return true
			} else if target.IsValid() {
				break
			}
		}
	}

	for _, off := range knightOffsets {
		to := Coord{from.File + off.f, from.Rank + off.r}
		if to.IsValid() && board.At(to).Color != side && board.At(to).Name == Knight {
			return true
		}
	}

	dir := 1
	if side == Black {
		dir = -1
	}
	pawnSquares := [...]Coord{{from.File - 1, from.Rank + dir}, {from.File + 1, from.Rank + dir}}
	for _, square := range pawnSquares {
		piece := board.At(square)
		if piece != nil && piece.Name == Pawn && piece.Color != side {
			return true
		}
	}

	return false
}
func (board *Board) InCheckmate() bool {
	return board.InCheck(board.SideToMove) && len(board.Moves()) == 0
}
func (board *Board) InStalemate() bool {
	return !board.InCheck(board.SideToMove) && len(board.Moves()) == 0
}
func (board *Board) GameOver() bool {
	count := len(board.Moves())
	return (board.InCheck(board.SideToMove) && count == 0) || (!board.InCheck(board.SideToMove) && count == 0)
}

func (board *Board) PseudoMoves(types ...PieceName) []Move {
	moveSet := make([]Move, 0, 128)

	pieces := board.pieceIndices(board.SideToMove, types...)

	moveSet = append(moveSet, board.pawnMoves(pieces)...)
	moveSet = append(moveSet, board.knightMoves(pieces)...)
	moveSet = append(moveSet, board.slidingMoves(pieces)...)
	moveSet = append(moveSet, board.kingMoves(pieces)...)

	return moveSet
}
func (board *Board) Moves() []Move {
	pseudoMoves := board.PseudoMoves()
	moveSet := make([]Move, 0, len(pseudoMoves))
	side := board.SideToMove

	for _, move := range pseudoMoves {
		isLegal := true

		if move.CastlesTo.IsValid() {
			off := 1
			if move.CastlesTo == Queenside {
				off = -1
			}

			if board.InCheck(side) {
				isLegal = false
			} else {
				if actual := board.MakeMove(Move{move.From, Coord{move.From.File + off, move.From.Rank}, MoveFlags{Moves: King}}); !actual.IsValid() {
					isLegal = false
				}

				if board.InCheck(side) {
					isLegal = false
				}
				board.UnmakeMove()
			}
		}

		var actual Move
		if isLegal {
			if actual = board.MakeMove(move); !actual.IsValid() {
				panic(fmt.Sprintf("invalid move generated: %v", actual))
			}

			if board.InCheck(side) {
				isLegal = false
			}
			if m := board.UnmakeMove(); !m.IsValid() {
				fmt.Println("Unable to unmake move")
			}
		}

		if isLegal {
			moveSet = append(moveSet, actual)
		}
	}

	return moveSet
}
func (board *Board) CountMoves(depth int) (int, []int) {
	if depth <= 0 {
		return 1, nil
	}

	moves := board.Moves()
	// if depth == 1 {
	// 	return len(moves), []int{}
	// }

	count, breakdown := 0, make([]int, 8)
	for _, move := range moves {
		if actual := board.MakeMove(move); !actual.IsValid() {
			panic(fmt.Sprintf("invalid move generated: %v", actual))
		}

		if move.Captures.IsValid() {
			breakdown[0]++
		}
		if move.IsEnPassant {
			breakdown[1]++
		}
		if move.CastlesTo.IsValid() {
			breakdown[2]++
		}
		if move.PromotesTo.IsValid() {
			breakdown[3]++
		}

		amount, parts := board.CountMoves(depth - 1)
		count += amount
		for i := 0; i < len(parts); i++ {
			breakdown[i] += parts[i]
		}

		board.UnmakeMove()
	}

	return count, breakdown
}

func (board *Board) updateState(move Move) {
	board.EnPassantTarget = Coord{0, 0}
	if diff := move.From.Rank - move.To.Rank; move.Moves == Pawn && diff/2 != 0 {
		left, right := board.At(Coord{move.To.File + 1, move.To.Rank}), board.At(Coord{move.To.File - 1, move.To.Rank})
		if (left != nil && left.Name == Pawn && left.Color != board.SideToMove) || (right != nil && right.Name == Pawn && right.Color != board.SideToMove) {
			board.EnPassantTarget = Coord{move.To.File, move.To.Rank + diff/2}
		}
	}

	if move.Moves == King {
		board.CastleRights.Disallow(board.SideToMove, Kingside)
		board.CastleRights.Disallow(board.SideToMove, Queenside)
	} else if move.Moves == Rook || move.Captures == Rook {
		var rook Coord
		if move.Moves == Rook {
			rook = move.From
		} else {
			rook = move.To
		}

		switch rook {
		case Coord{1, 1}:
			board.CastleRights.Disallow(White, Queenside)
		case Coord{8, 1}:
			board.CastleRights.Disallow(White, Kingside)
		case Coord{1, 8}:
			board.CastleRights.Disallow(Black, Queenside)
		case Coord{8, 8}:
			board.CastleRights.Disallow(Black, Kingside)
		}
	}

	board.SideToMove ^= 0b11
	if board.SideToMove == White {
		board.FullmoveCounter++
	}
	if move.Moves == Pawn || board.At(move.To).IsValid() {
		board.HalfmoveClock = 0
	} else {
		board.HalfmoveClock++
	}
}
func (board *Board) MakeMove(move Move) (actual Move) {
	if !move.IsValid() {
		return
	}

	movePiece := *board.At(move.From)
	actual.Moves = movePiece.Name
	if move.PromotesTo.IsValid() {
		movePiece = Piece{movePiece.Color, move.PromotesTo}
		actual.PromotesTo = movePiece.Name
	}

	capturePiece := *board.At(move.To)
	actual.Captures = capturePiece.Name

	if move.CastlesTo.IsValid() {
		var rookFrom, rookTo Coord
		if move.CastlesTo == Kingside {
			rookFrom = Coord{8, move.To.Rank}
		} else {
			rookFrom = Coord{1, move.To.Rank}
		}

		off := int(move.CastlesTo)*2 - 3
		rookTo = Coord{move.To.File + off, move.To.Rank}
		if !rookFrom.IsValid() || !rookTo.IsValid() {
			return
		}

		*board.At(rookTo) = *board.At(rookFrom)
		*board.At(rookFrom) = Piece{0, 0}
		actual.CastlesTo = move.CastlesTo
	} else if movePiece.Name == Pawn && move.To == board.EnPassantTarget {
		*board.At(Coord{move.To.File, move.From.Rank}) = Piece{0, 0}
		actual.IsEnPassant = true
	}

	*board.At(move.To) = movePiece
	*board.At(move.From) = Piece{0, 0}

	actual.To = move.To
	actual.From = move.From
	actual.OffersDraw = move.OffersDraw

	board.history = append(board.history, BoardState{actual, board.BoardData})
	board.updateState(actual)

	return
}
func (board *Board) UnmakeMove() Move {
	if len(board.history) == 0 {
		return Move{}
	}

	i := len(board.history) - 1
	state := board.history[i]
	if !state.Move.IsValid() {
		return Move{}
	}

	to, from := board.At(state.To), board.At(state.From)
	if state.PromotesTo.IsValid() {
		*from = Piece{to.Color, Pawn}
	} else {
		*from = *to
	}

	if state.Captures.IsValid() {
		*to = Piece{board.SideToMove, state.Captures}
	} else {
		*to = Piece{0, 0}
	}

	if state.CastlesTo == Kingside {
		rook := board.At(Coord{6, state.To.Rank})
		corner := board.At(Coord{8, state.To.Rank})
		*corner = *rook
		*rook = Piece{0, 0}
	} else if state.CastlesTo == Queenside {
		rook := board.At(Coord{4, state.To.Rank})
		corner := board.At(Coord{1, state.To.Rank})
		*corner = *rook
		*rook = Piece{0, 0}
	}

	if state.IsEnPassant {
		passant := Coord{state.To.File, state.From.Rank}
		*board.At(passant) = Piece{board.SideToMove, Pawn}
	}

	board.BoardData = state.BoardData
	board.history = board.history[:i]

	return state.Move
}
