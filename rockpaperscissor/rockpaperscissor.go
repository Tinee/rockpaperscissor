package rockpaperscissor

import (
	"time"
)

type Result struct {
	Winner   Player
	Loser    Player
	Draw     bool
	PlayedAt time.Time
}

type Mover interface {
	Move() Move
}

type Player interface {
	Name() string
	Mover
}

// Move is a weapon type of int
type Move int

// Different weapons that the game has.
const (
	Rock Move = iota
	Paper
	Scissor
)

func Play(p1, p2 Player) Result {
	return plays(p1, p2)
}

func plays(p1, p2 Player) Result {
	if p1.Move() == p2.Move() {
		return Result{Draw: true}
	}

	if beats(p2) == p1.Move() {
		return Result{
			Winner: p2,
			Loser:  p1,
			Draw:   false,
		}
	}

	return Result{
		Winner: p1,
		Loser:  p2,
		Draw:   false,
	}
}

func beats(m Mover) Move {
	switch m.Move() {
	case Rock:
		return Scissor
	case Scissor:
		return Paper
	case Paper:
		return Rock
	default:
		// This sucks? This will never happen!
		return Rock
	}
}
