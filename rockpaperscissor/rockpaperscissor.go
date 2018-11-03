package rockpaperscissor

// Mover is someone that have moves.
type Mover interface {
	Move() Move
}

// Decider is someone that is smart enough to make their own
// decisions when they either lose, wins or draws/
type Decider interface {
	Decide(Outcome)
}

// Player is someone that can play Rock Paper Scissor,
// He can also determine how to behave when he gets the result.
type Player interface {
	Mover
	Decider
}

// Play takes players and determine who the winner is.
func Play(p1, p2 Player) {
	if p1.Move() == p2.Move() {
		p1.Decide(Draw)
		p2.Decide(Draw)
	}

	if beats(p2) == p1.Move() {
		p1.Decide(Lost)
		p2.Decide(Won)
	}

	p1.Decide(Won)
	p2.Decide(Lost)
}

// Move is a move type of int
type Move int

// Different moves that the game has.
const (
	Rock Move = iota
	Paper
	Scissor
)

// Outcome is a the different outcomes that can happen when people play.
type Outcome int

// different outcomes.
const (
	Won Outcome = iota
	Lost
	Draw
)

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
