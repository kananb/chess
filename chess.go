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

		start := coordFromIndex(i)
		end := NewCoord(start.File(), start.Rank()+dir)
		if end != NoCoord && *board.At(end) == NoPiece {
			moveSet = append(moveSet, NewMove(start, end))

			end = NewCoord(start.File(), start.Rank()+dir*2)
			canMoveDouble := (piece.Color() == White && start.Rank() == 1) || (piece.Color() == Black && start.Rank() == 6)
			if end != NoCoord && canMoveDouble && *board.At(end) == NoPiece {
				moveSet = append(moveSet, NewMove(start, end))
			}
		}
		for off := 3; off > 0; off -= 2 {
			end := NewCoord(start.File()+off-2, start.Rank()+dir)
			if end != NoCoord && *board.At(end) != NoPiece && board.At(end).Color() != piece.Color() || board.EnPassantTarget == end {
				moveSet = append(moveSet, NewMove(start, end))
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

		start := coordFromIndex(i)
		offsets := [...]struct{ f, r int }{{1, 2}, {2, 1}, {2, -1}, {1, -2}, {-1, -2}, {-2, -1}, {-2, 1}, {-1, 2}}
		for _, off := range offsets {
			end := NewCoord(start.File()+off.f, start.Rank()+off.r)
			if end != NoCoord && board.At(end).Color() != piece.Color() {
				moveSet = append(moveSet, NewMove(start, end))
			}
		}
	}

	return moveSet
}
func (board *Board) genSlidingMoves(pieces []int) []Move {
	moveSet := make([]Move, 0, 128)
	slideDirections := [...]struct{ f, r int }{
		{1, 1}, {1, -1}, {-1, -1}, {-1, 1},
		{0, 1}, {1, 0}, {0, -1}, {-1, 0},
	}

	for _, i := range pieces {
		piece := board.squares[i]
		if piece.Type() != Bishop && piece.Type() != Rook && piece.Type() != Queen {
			continue
		}

		var directions []struct{ f, r int }
		switch piece.Type() {
		case Bishop:
			directions = slideDirections[:4]
		case Rook:
			directions = slideDirections[4:]
		default:
			directions = slideDirections[:]
		}

		start := coordFromIndex(i)
		for _, dir := range directions {
			for off := 1; ; off++ {
				end := NewCoord(start.File()+dir.f*off, start.Rank()+dir.r*off)
				if end == NoCoord || board.At(end).Color() == piece.Color() {
					break
				}

				moveSet = append(moveSet, NewMove(start, end))
				if *board.At(end) != NoPiece {
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

		start := coordFromIndex(i)
		for x := -1; x < 2; x++ {
			for y := -1; y < 2; y++ {
				end := NewCoord(start.File()+x, start.Rank()+y)
				if end != NoCoord && board.At(end).Color() != piece.Color() {
					moveSet = append(moveSet, NewMove(start, end))
				}
			}
		}

		break
	}

	return moveSet
}
func (board *Board) GenMoves() []Move {
	moveSet := make([]Move, 0, 128)

	pieces := make([]int, 0, 16)
	for i := 0; i < len(board.squares); i++ {
		if board.squares[i].Color() == board.SideToMove {
			pieces = append(pieces, i)
		}
	}

	moveSet = append(moveSet, board.genPawnMoves(pieces)...)
	moveSet = append(moveSet, board.genKnightMoves(pieces)...)
	moveSet = append(moveSet, board.genSlidingMoves(pieces)...)
	moveSet = append(moveSet, board.genKingMoves(pieces)...)
	return moveSet
}

func (board *Board) MakeMove(move Move) bool {
	piece := board.At(move.Start)
	capture := *board.At(move.End) != NoPiece
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
