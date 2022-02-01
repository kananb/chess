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

func (board *Board) getPieceIndices(side SideColor, types ...PieceName) []int {
	pieces := make([]int, 0, 16)

	for i := 0; i < len(board.squares); i++ {
		if board.squares[i].Color != side {
			continue
		}

		if len(types) == 0 {
			pieces = append(pieces, i)
			continue
		}
		for _, t := range types {
			if board.squares[i].Name == t {
				pieces = append(pieces, i)
				break
			}
		}
	}

	return pieces
}
func (board *Board) genPawnMoves(pieces []int) []Move {
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
		if to.IsValid() && board.At(to).IsValid() {
			if to.Rank == 0 || to.Rank == 7 {
				for _, t := range promoteTypes {
					moveSet = append(moveSet, Move{from, to, MoveFlags{Moves: Pawn, Promotes: t}})
				}
			} else {
				moveSet = append(moveSet, Move{from, to, MoveFlags{Moves: Pawn}})
			}

			to = Coord{from.File, from.Rank + dir*2}
			canMoveDouble := (piece.Color == White && from.Rank == 1) || (piece.Color == Black && from.Rank == 6)
			if to.IsValid() && canMoveDouble && board.At(to).IsValid() {
				moveSet = append(moveSet, Move{from, to, MoveFlags{Moves: Pawn}})
			}
		}
		for off := 3; off > 0; off -= 2 {
			to := Coord{from.File + off - 2, from.Rank + dir}
			if to.IsValid() && (board.At(to).IsValid() && board.At(to).Color != piece.Color || board.EnPassantTarget == to) {
				if to.Rank == 0 || to.Rank == 7 {
					for _, t := range promoteTypes {
						moveSet = append(moveSet, Move{from, to, MoveFlags{Moves: Pawn, Captures: board.At(to).Name, Promotes: t}})
					}
				} else {
					enPassant := false
					if board.EnPassantTarget == to {
						enPassant = true
					}
					moveSet = append(moveSet, Move{from, to, MoveFlags{Moves: Pawn, Captures: board.At(to).Name, EnPassant: enPassant}})
				}
			}
		}
	}

	return moveSet
}
func (board *Board) genKnightMoves(pieces []int) []Move {
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
func (board *Board) genSlidingMoves(pieces []int) []Move {
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
func (board *Board) genKingMoves(pieces []int) []Move {
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
			if board.CastleRights.Can(piece.Color, side) && !board.At(between).IsValid() && board.At(to).IsValid() {
				moveSet = append(moveSet, Move{from, to, MoveFlags{Moves: King, Castle: side}})
			}
		}
		checkCastle(Kingside, 1)
		checkCastle(Queenside, -1)

		break
	}

	return moveSet
}
func (board *Board) genPseudoMoves(types ...PieceName) []Move {
	moveSet := make([]Move, 0, 128)

	pieces := board.getPieceIndices(board.SideToMove, types...)

	moveSet = append(moveSet, board.genPawnMoves(pieces)...)
	moveSet = append(moveSet, board.genKnightMoves(pieces)...)
	moveSet = append(moveSet, board.genSlidingMoves(pieces)...)
	moveSet = append(moveSet, board.genKingMoves(pieces)...)

	return moveSet
}
func (board *Board) kingInCheck(side SideColor) bool {
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
func (board *Board) GenMoves() []Move {
	pseudoMoves := board.genPseudoMoves()
	moveSet := make([]Move, 0, len(pseudoMoves))
	side := board.SideToMove

	for _, move := range pseudoMoves {
		isLegal := true

		if move.Castle.IsValid() {
			off := 1
			if move.Castle == Queenside {
				off = -1
			}

			if board.kingInCheck(side) {
				isLegal = false
			} else {
				if err := board.MakeMove(Move{move.From, Coord{move.From.File + off, move.From.Rank}, MoveFlags{Moves: King}}); err != nil {
					isLegal = false
				}

				if board.kingInCheck(side) {
					isLegal = false
				}
				board.UnmakeMove()
			}
		}

		if isLegal {
			if err := board.MakeMove(move); err != nil {
				panic(fmt.Sprintf("invalid move generated: %v", err))
			}

			if board.kingInCheck(side) {
				isLegal = false
			}
			if _, ok := board.UnmakeMove(); !ok {
				fmt.Println("Unable to unmake move")
			}
		}

		if isLegal {
			moveSet = append(moveSet, move)
		}
	}

	return moveSet
}
func (board *Board) CountMoves(depth int) (int, []int) {
	if depth <= 0 {
		return 1, nil
	}

	moves := board.GenMoves()
	// if depth == 1 {
	// 	return len(moves), []int{}
	// }

	count, breakdown := 0, make([]int, 8)
	for _, move := range moves {
		if err := board.MakeMove(move); err != nil {
			panic("invalid move generated")
		}

		if move.Captures.IsValid() {
			breakdown[0]++
		}
		if move.EnPassant {
			breakdown[1]++
		}
		if move.Castle.IsValid() {
			breakdown[2]++
		}
		if move.Promotes.IsValid() {
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

func (board *Board) MakeMove(move Move) (err error) {
	if !move.IsValid() {
		return fmt.Errorf("invalid move coordinates")
	}
	turn := moveState{move, board.CastleRights, board.EnPassantTarget, board.HalfmoveClock}

	movePiece := *board.At(move.From)
	if move.Moves != movePiece.Name {
		return fmt.Errorf("piece types don't match")
	}

	if diff := move.From.File - move.To.File; move.Castle.IsValid() {
		if !board.CastleRights.Can(board.SideToMove, move.Castle) {
			return fmt.Errorf("castling on that side is not allowed: %v [%v]", move.Castle, board.String())
		} else if move.Captures.IsValid() {
			return fmt.Errorf("cannot castle and capture")
		} else if move.Promotes.IsValid() {
			return fmt.Errorf("cannot castle and promote")
		} else if move.EnPassant {
			return fmt.Errorf("cannot castle and en passant")
		} else if movePiece.Name != King || diff/2 == 0 {
			return fmt.Errorf("invalid castle conditions")
		}
	} else if movePiece.Name == King && diff/2 != 0 {
		return fmt.Errorf("cannot castle implicitly")
	}

	capturePiece := *board.At(move.To)
	if move.Captures != capturePiece.Name {
		return fmt.Errorf("capture types don't match")
	}

	if promoteType := move.Promotes; promoteType.IsValid() {
		if movePiece.Name != Pawn || (move.To.Rank != 7 && move.To.Rank != 0) || promoteType == King || promoteType == Pawn {
			return fmt.Errorf("invalid promotion: %v [%v] on rank %d", movePiece, promoteType, move.To.Rank)
		}

		movePiece = Piece{movePiece.Color, promoteType}
	} else if movePiece.Name == Pawn && (move.To.Rank == 7 || move.To.Rank == 0) {
		return fmt.Errorf("pawns on the last rank must promote")
	}

	if move.EnPassant && move.To != board.EnPassantTarget {
		return fmt.Errorf("cannot en passant in this position: %v", board.String())
	} else if move.EnPassant && movePiece.Name != Pawn {
		return fmt.Errorf("cannot en passant non-pawn piece")
	}

	*board.At(move.To) = movePiece
	*board.At(move.From) = Piece{0, 0}

	if move.Castle == Kingside {
		rook := board.At(Coord{7, move.To.Rank})
		home := board.At(Coord{5, move.To.Rank})
		*home = *rook
		*rook = Piece{0, 0}
	} else if move.Castle == Queenside {
		rook := board.At(Coord{0, move.To.Rank})
		home := board.At(Coord{3, move.To.Rank})
		*home = *rook
		*rook = Piece{0, 0}
	} else if movePiece.Name == Pawn && move.To == board.EnPassantTarget {
		*board.At(Coord{move.To.File, move.From.Rank}) = Piece{0, 0}
	}

	board.EnPassantTarget = Coord{0, 0}
	if diff := move.From.Rank - move.To.Rank; movePiece.Name == Pawn && diff/2 != 0 {
		left, right := board.At(Coord{move.To.File + 1, move.To.Rank}), board.At(Coord{move.To.File - 1, move.To.Rank})
		if (left != nil && left.Name == Pawn && left.Color != movePiece.Color) || (right != nil && right.Name == Pawn && right.Color != movePiece.Color) {
			board.EnPassantTarget = Coord{move.To.File, move.To.Rank + diff/2}
		}
	}

	if movePiece.Name == King {
		board.CastleRights.Disallow(board.SideToMove, Kingside)
		board.CastleRights.Disallow(board.SideToMove, Queenside)
	} else if movePiece.Name == Rook {
		switch move.From {
		case Coord{0, 0}:
			board.CastleRights.Disallow(White, Queenside)
		case Coord{7, 0}:
			board.CastleRights.Disallow(White, Kingside)
		case Coord{0, 7}:
			board.CastleRights.Disallow(Black, Queenside)
		case Coord{7, 7}:
			board.CastleRights.Disallow(Black, Kingside)
		}
	} else if move.Captures == Rook {
		switch move.To {
		case Coord{0, 0}:
			board.CastleRights.Disallow(White, Queenside)
		case Coord{7, 0}:
			board.CastleRights.Disallow(White, Kingside)
		case Coord{0, 7}:
			board.CastleRights.Disallow(Black, Queenside)
		case Coord{7, 7}:
			board.CastleRights.Disallow(Black, Kingside)
		}
	}

	board.SideToMove ^= 0b11
	if movePiece.Color == Black {
		board.FullmoveCounter++
	}
	if movePiece.Name == Pawn || board.At(move.To).IsValid() {
		board.HalfmoveClock = 0
	} else {
		board.HalfmoveClock++
	}

	if board.history == nil {
		board.history = make([]moveState, 0, 128)
	}
	board.history = append(board.history, turn)

	return
}
func (board *Board) UnmakeMove() (Move, bool) {
	if len(board.history) == 0 {
		return Move{}, false
	}

	i := len(board.history) - 1
	turn := board.history[i]

	to, from := board.At(turn.To), board.At(turn.From)
	if !turn.Promotes.IsValid() {
		*from = *to
	} else {
		*from = Piece{to.Color, Pawn}
	}

	if turn.Captures.IsValid() {
		*to = Piece{board.SideToMove, turn.Captures}
	} else {
		*to = Piece{0, 0}
	}

	if turn.Castle == Kingside {
		rook := board.At(Coord{5, turn.To.Rank})
		corner := board.At(Coord{7, turn.To.Rank})
		*corner = *rook
		*rook = Piece{0, 0}
	} else if turn.Castle == Queenside {
		rook := board.At(Coord{3, turn.To.Rank})
		corner := board.At(Coord{0, turn.To.Rank})
		*corner = *rook
		*rook = Piece{0, 0}
	}

	if turn.EnPassant {
		passant := Coord{turn.To.File, turn.From.Rank}
		*board.At(passant) = Piece{board.SideToMove, Pawn}
		board.EnPassantTarget = Coord{turn.To.File, turn.From.Rank}
	}

	board.SideToMove ^= 0b11
	if board.SideToMove == Black {
		board.FullmoveCounter--
	}

	board.CastleRights = turn.CastleRights
	board.EnPassantTarget = turn.EnPassantTarget
	board.HalfmoveClock = turn.HalfmoveClock

	board.history = board.history[:i]
	return turn.Move, true
}
