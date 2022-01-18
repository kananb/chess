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
var promoteTypes = [...]PieceType{
	Knight, Bishop, Rook, Queen,
}

func (board *Board) getPieceIndices(types ...PieceType) []int {
	pieces := make([]int, 0, 16)

	for i := 0; i < len(board.squares); i++ {
		if board.squares[i].Color() == board.SideToMove {
			if len(types) == 0 {
				pieces = append(pieces, i)
				continue
			}
			for _, t := range types {
				if board.squares[i].Type() == t {
					pieces = append(pieces, i)
					break
				}
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
		if piece.Type() != Pawn {
			continue
		}

		from := coordFromIndex(i)
		to := NewCoord(from.File(), from.Rank()+dir)
		if to != NoCoord && *board.At(to) == NoPiece {
			if to.Rank() == 0 || to.Rank() == 7 {
				for _, t := range promoteTypes {
					moveSet = append(moveSet, NewMove(from, to, MoveFlags{Moves: Pawn, Promotes: t}))
				}
			} else {
				moveSet = append(moveSet, NewMove(from, to, MoveFlags{Moves: Pawn}))
			}

			to = NewCoord(from.File(), from.Rank()+dir*2)
			canMoveDouble := (piece.Color() == White && from.Rank() == 1) || (piece.Color() == Black && from.Rank() == 6)
			if to != NoCoord && canMoveDouble && *board.At(to) == NoPiece {
				moveSet = append(moveSet, NewMove(from, to, MoveFlags{Moves: Pawn}))
			}
		}
		for off := 3; off > 0; off -= 2 {
			to := NewCoord(from.File()+off-2, from.Rank()+dir)
			if to != NoCoord && (*board.At(to) != NoPiece && board.At(to).Color() != piece.Color() || board.EnPassantTarget == to) {
				if to.Rank() == 0 || to.Rank() == 7 {
					for _, t := range promoteTypes {
						moveSet = append(moveSet, NewMove(from, to, MoveFlags{Moves: Pawn, Captures: board.At(to).Type(), Promotes: t}))
					}
				} else {
					enPassant := false
					if board.EnPassantTarget == to {
						enPassant = true
					}
					moveSet = append(moveSet, NewMove(from, to, MoveFlags{Moves: Pawn, Captures: board.At(to).Type(), EnPassant: enPassant}))
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
		if piece.Type() != Knight {
			continue
		}

		from := coordFromIndex(i)
		for _, off := range knightOffsets {
			to := NewCoord(from.File()+off.f, from.Rank()+off.r)
			if to != NoCoord && board.At(to).Color() != piece.Color() {
				moveSet = append(moveSet, NewMove(from, to, MoveFlags{Moves: Knight, Captures: board.At(to).Type()}))
			}
		}
	}

	return moveSet
}
func (board *Board) genSlidingMoves(pieces []int) []Move {
	moveSet := make([]Move, 0, 128)

	for _, i := range pieces {
		piece := board.squares[i]
		if piece.Type() != Bishop && piece.Type() != Rook && piece.Type() != Queen {
			continue
		}

		di, df := 0, 8
		if piece.Type() == Bishop {
			df = 4
		} else if piece.Type() == Rook {
			di = 4
		}

		from := coordFromIndex(i)
		for d := di; d < df; d++ {
			for off := 1; ; off++ {
				to := NewCoord(from.File()+slideDirections[d].f*off, from.Rank()+slideDirections[d].r*off)
				if to == NoCoord || board.At(to).Color() == piece.Color() {
					break
				}

				moveSet = append(moveSet, NewMove(from, to, MoveFlags{Moves: piece.Type(), Captures: board.At(to).Type()}))
				if *board.At(to) != NoPiece {
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
		if piece.Type() != King {
			continue
		}

		from := coordFromIndex(i)
		for x := -1; x < 2; x++ {
			for y := -1; y < 2; y++ {
				to := NewCoord(from.File()+x, from.Rank()+y)
				if to != NoCoord && board.At(to).Color() != piece.Color() {
					moveSet = append(moveSet, NewMove(from, to, MoveFlags{Moves: King, Captures: board.At(to).Type()}))
				}
			}
		}

		checkCastle := func(side CastleSide, dir int) {
			between, to := NewCoord(from.File()+1*dir, from.Rank()), NewCoord(from.File()+2*dir, from.Rank())
			if board.CastleRights.Can(piece.Color(), side) && *board.At(between) == NoPiece && *board.At(to) == NoPiece {
				moveSet = append(moveSet, NewMove(from, to, MoveFlags{Moves: King, Castle: side}))
			}
		}
		checkCastle(Kingside, 1)
		checkCastle(Queenside, -1)

		break
	}

	return moveSet
}
func (board *Board) genPseudoMoves(types ...PieceType) []Move {
	moveSet := make([]Move, 0, 128)

	pieces := board.getPieceIndices(types...)

	moveSet = append(moveSet, board.genPawnMoves(pieces)...)
	moveSet = append(moveSet, board.genKnightMoves(pieces)...)
	moveSet = append(moveSet, board.genSlidingMoves(pieces)...)
	moveSet = append(moveSet, board.genKingMoves(pieces)...)

	return moveSet
}
func (board *Board) kingInCheck(side SideColor) bool {
	var from Coord
	for i := 0; i < len(board.squares); i++ {
		if board.squares[i].Color() == side && board.squares[i].Type() == King {
			from = coordFromIndex(i)
			break
		}
	}

	for i, dir := range slideDirections {
		for off := 1; ; off++ {
			to := NewCoord(from.File()+dir.f*off, from.Rank()+dir.r*off)
			if to == NoCoord || board.At(to).Color() == side {
				break
			}

			target := board.At(to).Type()
			if target == Queen || (i < 4 && target == Bishop) || (i >= 4 && target == Rook) {
				return true
			} else if target != NoType {
				break
			}
		}
	}

	for _, off := range knightOffsets {
		to := NewCoord(from.File()+off.f, from.Rank()+off.r)
		if to != NoCoord && board.At(to).Color() != side && board.At(to).Type() == Knight {
			return true
		}
	}

	dir := 1
	if side == Black {
		dir = -1
	}
	pawnSquares := [...]Coord{NewCoord(from.File()-1, from.Rank()+dir), NewCoord(from.File()+1, from.Rank()+dir)}
	for _, square := range pawnSquares {
		piece := board.At(square)
		if piece != nil && piece.Type() == Pawn && piece.Color() != side {
			return true
		}
	}

	return false
}
func (board *Board) GenMoves() []Move {
	pseudoMoves := board.genPseudoMoves()
	moveSet := make([]Move, 0, len(pseudoMoves))

	for _, move := range pseudoMoves {
		isLegal := true

		if move.CastlesTo() != NoCastle {
			off := 1
			if move.CastlesTo() == Queenside {
				off = -1
			}

			if board.kingInCheck(board.SideToMove) {
				isLegal = false
			} else {
				if err := board.MakeMove(NewMove(move.From, NewCoord(move.From.File()+off, move.From.Rank()), MoveFlags{Moves: King})); err != nil {
					isLegal = false
				}

				if board.kingInCheck(board.SideToMove) {
					isLegal = false
				}
				board.UnmakeMove()
			}
		}

		if isLegal {
			if err := board.MakeMove(move); err != nil {
				panic(fmt.Sprintf("invalid move generated: %v", err))
			}

			if board.kingInCheck(board.SideToMove) {
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
func (board *Board) CountMoves(depth, maxDepth int) (int, []int) {
	if depth >= maxDepth {
		return 0, nil
	}

	moves := board.GenMoves()
	count, breakdown := len(moves), make([]int, 8)
	for _, move := range moves {
		if err := board.MakeMove(move); err != nil {
			panic("invalid move generated")
		}

		if move.Captures() != NoType {
			breakdown[0]++
		}
		if move.IsEnPassant() {
			breakdown[1]++
		}
		if move.CastlesTo() != NoCastle {
			breakdown[2]++
		}
		if move.PromotesTo() != NoType {
			breakdown[3]++
		}

		amount, parts := board.CountMoves(depth+1, maxDepth)
		count += amount
		for i := 0; i < len(parts); i++ {
			breakdown[i] += parts[i]
		}

		board.UnmakeMove()
	}

	return count, breakdown
}

func (board *Board) MakeMove(move Move) (err error) {
	if move.From == NoCoord || move.To == NoCoord {
		return fmt.Errorf("invalid move coordinates")
	}

	movePiece := *board.At(move.From)
	if move.Moves() != movePiece.Type() {
		return fmt.Errorf("piece types don't match")
	}

	if diff := move.From.File() - move.To.File(); move.CastlesTo() != NoCastle {
		if !board.CastleRights.Can(board.SideToMove, move.CastlesTo()) {
			return fmt.Errorf("castling on that side is not allowed: %v [%v]", move.CastlesTo(), board.FEN())
		} else if move.Captures() != NoType {
			return fmt.Errorf("cannot castle and capture")
		} else if move.PromotesTo() != NoType {
			return fmt.Errorf("cannot castle and promote")
		} else if move.IsEnPassant() {
			return fmt.Errorf("cannot castle and en passant")
		} else if movePiece.Type() != King || diff/2 == 0 {
			return fmt.Errorf("invalid castle conditions")
		}
	} else if movePiece.Type() == King && diff/2 != 0 {
		return fmt.Errorf("cannot castle implicitly")
	}

	capturePiece := *board.At(move.To)
	if move.Captures() != capturePiece.Type() {
		return fmt.Errorf("capture types don't match")
	}

	if promoteType := move.PromotesTo(); promoteType != NoType {
		if movePiece.Type() != Pawn || (move.To.Rank() != 7 && move.To.Rank() != 0) || promoteType == King || promoteType == Pawn {
			return fmt.Errorf("invalid promotion: %v [%v] on rank %d", movePiece, promoteType, move.To.Rank())
		}

		movePiece = NewPiece(movePiece.Color(), promoteType)
	} else if movePiece.Type() == Pawn && (move.To.Rank() == 7 || move.To.Rank() == 0) {
		return fmt.Errorf("pawns on the last rank must promote")
	}

	if move.IsEnPassant() && move.To != board.EnPassantTarget {
		return fmt.Errorf("cannot en passant in this position: %v", board.FEN())
	} else if move.IsEnPassant() && movePiece.Type() != Pawn {
		return fmt.Errorf("cannot en passant non-pawn piece")
	}

	*board.At(move.To) = movePiece
	*board.At(move.From) = NoPiece

	if move.CastlesTo() == Kingside {
		rook := board.At(NewCoord(7, move.To.Rank()))
		home := board.At(NewCoord(5, move.To.Rank()))
		*home = *rook
		*rook = NoPiece
	} else if move.CastlesTo() == Queenside {
		rook := board.At(NewCoord(0, move.To.Rank()))
		home := board.At(NewCoord(3, move.To.Rank()))
		*home = *rook
		*rook = NoPiece
	} else if movePiece.Type() == Pawn && move.To == board.EnPassantTarget {
		*board.At(NewCoord(move.To.File(), move.From.Rank())) = NoPiece
	}

	board.EnPassantTarget = NoCoord
	if diff := move.From.Rank() - move.To.Rank(); movePiece.Type() == Pawn && diff/2 != 0 {
		left, right := board.At(NewCoord(move.To.File()+1, move.To.Rank())), board.At(NewCoord(move.To.File()-1, move.To.Rank()))
		if (left != nil && left.Type() == Pawn && left.Color() != movePiece.Color()) || (right != nil && right.Type() == Pawn && right.Color() != movePiece.Color()) {
			board.EnPassantTarget = NewCoord(move.To.File(), move.To.Rank()+diff/2)
		}
	}

	move.flags = (move.flags &^ (prevRightsMask << prevRightsShift)) | (uint32(board.CastleRights.ToWord()) << prevRightsShift)
	if movePiece.Type() == King {
		board.CastleRights.Disallow(board.SideToMove, Kingside)
		board.CastleRights.Disallow(board.SideToMove, Queenside)
	} else if movePiece.Type() == Rook {
		switch move.From {
		case NewCoord(0, 0):
			board.CastleRights.Disallow(White, Queenside)
		case NewCoord(7, 0):
			board.CastleRights.Disallow(White, Kingside)
		case NewCoord(0, 7):
			board.CastleRights.Disallow(Black, Queenside)
		case NewCoord(7, 7):
			board.CastleRights.Disallow(Black, Kingside)
		}
	} else if move.Captures() == Rook {
		switch move.To {
		case NewCoord(0, 0):
			board.CastleRights.Disallow(White, Queenside)
		case NewCoord(7, 0):
			board.CastleRights.Disallow(White, Kingside)
		case NewCoord(0, 7):
			board.CastleRights.Disallow(Black, Queenside)
		case NewCoord(7, 7):
			board.CastleRights.Disallow(Black, Kingside)
		}
	}

	board.SideToMove ^= 0b11
	if movePiece.Color() == Black {
		board.FullmoveCounter++
	}
	if movePiece.Type() == Pawn || *board.At(move.To) != NoPiece {
		board.HalfmoveClock = 0
	} else {
		board.HalfmoveClock++
	}
	move.flags = (move.flags &^ (plyMask << plyShift)) | (uint32(board.HalfmoveClock) << plyShift)

	if board.history == nil {
		board.history = make([]Move, 0, 128)
	}
	board.history = append(board.history, move)

	return
}
func (board *Board) UnmakeMove() (Move, bool) {
	if len(board.history) == 0 {
		return Move{}, false
	}

	i := len(board.history) - 1
	move := board.history[i]

	to, from := board.At(move.To), board.At(move.From)
	if move.PromotesTo() == NoType {
		*from = *to
	} else {
		*from = NewPiece(to.Color(), Pawn)
	}

	if move.Captures() != NoType {
		*to = NewPiece(board.SideToMove, move.Captures())
	} else {
		*to = NoPiece
	}

	if move.CastlesTo() == Kingside {
		rook := board.At(NewCoord(5, move.To.Rank()))
		corner := board.At(NewCoord(7, move.To.Rank()))
		*corner = *rook
		*rook = NoPiece
	} else if move.CastlesTo() == Queenside {
		rook := board.At(NewCoord(3, move.To.Rank()))
		corner := board.At(NewCoord(0, move.To.Rank()))
		*corner = *rook
		*rook = NoPiece
	}

	board.SideToMove ^= 0b11
	if board.SideToMove == Black {
		board.FullmoveCounter--
	}

	if move.IsEnPassant() {
		passant := NewCoord(move.To.File(), move.From.Rank())
		*board.At(passant) = NewPiece(board.SideToMove, Pawn)
		board.EnPassantTarget = NewCoord(move.To.File(), move.From.Rank())
	}

	board.CastleRights = CastlesFromWord(uint8(move.flags >> prevRightsShift & prevRightsMask))
	board.HalfmoveClock = int(move.flags >> plyShift & plyMask)

	board.history = board.history[:i]
	return move, true
}
