package chess

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

		start := NewCoord(File(i%8), Rank(i/8))
		if board.squares[i+dir*8].Type() == NoPiece {
			moveSet = append(moveSet, NewMove(start, NewCoord(start.File(), Rank(int(start.Rank())+dir))))
			if (piece.Color() == White && start.Rank() == 1) || (piece.Color() == Black && start.Rank() == 6) && board.squares[i+dir*16].Type() == NoPiece {
				moveSet = append(moveSet, NewMove(start, NewCoord(start.File(), Rank(int(start.Rank())+dir*2))))
			}
		}
		for off := 3; off > 0; off -= 2 {
			end := NewCoord(File(int(start.File())+off-2), Rank(int(start.Rank())+dir))
			if board.At(end).Type() != NoPiece && board.At(end).Color() != piece.Color() || board.EnPassantTarget == end {
				moveSet = append(moveSet, NewMove(start, end))
			}
		}
	}

	return moveSet
}
func (board *Board) genHorseMoves(pieces []int) []Move {

	return nil
}
func (board *Board) genSlidingMoves(pieces []int) []Move {

	return nil
}
func (board *Board) genKingMoves(pieces []int) []Move {

	return nil
}
func (board *Board) GenMoves() []Move {
	moveSet := make([]Move, 0, 256)

	pieces := make([]int, 0, 16)
	for i := 0; i < len(board.squares); i++ {
		if board.squares[i].Color() == board.SideToMove {
			pieces = append(pieces, i)
		}
	}

	moveSet = append(moveSet, board.genPawnMoves(pieces)...)
	moveSet = append(moveSet, board.genHorseMoves(pieces)...)
	moveSet = append(moveSet, board.genSlidingMoves(pieces)...)
	moveSet = append(moveSet, board.genKingMoves(pieces)...)
	return moveSet
}

func (board *Board) MakeMove(move Move) bool {
	piece := board.At(move.Start)
	capture := board.At(move.End).Type() != NoPiece
	pawnMove := piece.Type() == Pawn

	*board.At(move.End) = *piece
	*board.At(move.Start) = Piece(0)

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
