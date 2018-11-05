package rockpaperscissor_test

import (
	"rockpaperscissor/rockpaperscissor"
	"testing"
)

type MockPlayer struct {
	name   string
	move   rockpaperscissor.Move
	output rockpaperscissor.Outcome
}

func (m *MockPlayer) Move() rockpaperscissor.Move       { return m.move }
func (m *MockPlayer) Name() string                      { return m.name }
func (m *MockPlayer) Decide(o rockpaperscissor.Outcome) { m.output = o }

func TestWeapon_Beats(t *testing.T) {
	tests := []struct {
		name                        string
		firstPlayerMove             rockpaperscissor.Move
		secondPlayerMove            rockpaperscissor.Move
		expectedFirstPlayerOutcome  rockpaperscissor.Outcome
		expectedSecondPlayerOutcome rockpaperscissor.Outcome
	}{
		{
			"Paper beats rock",
			rockpaperscissor.Paper,
			rockpaperscissor.Rock,
			rockpaperscissor.Won,
			rockpaperscissor.Lost,
		},
		{
			"Rock beats scissors",
			rockpaperscissor.Rock,
			rockpaperscissor.Scissor,
			rockpaperscissor.Won,
			rockpaperscissor.Lost,
		},
		{
			"Scissor beats paper",
			rockpaperscissor.Scissor,
			rockpaperscissor.Paper,
			rockpaperscissor.Won,
			rockpaperscissor.Lost,
		},
		{
			"Paper lose against scissors",
			rockpaperscissor.Paper,
			rockpaperscissor.Scissor,
			rockpaperscissor.Lost,
			rockpaperscissor.Won,
		},
		{
			"Rock lose against paper",
			rockpaperscissor.Rock,
			rockpaperscissor.Paper,
			rockpaperscissor.Lost,
			rockpaperscissor.Won,
		},
		{
			"Same weapon, it's a draw.",
			rockpaperscissor.Rock,
			rockpaperscissor.Rock,
			rockpaperscissor.Draw,
			rockpaperscissor.Draw,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			firstPlayer := MockPlayer{
				name: "First Player Name",
				move: tt.firstPlayerMove,
			}
			secondPlayer := MockPlayer{
				name: "First Player Name",
				move: tt.secondPlayerMove,
			}
			rockpaperscissor.Play(&firstPlayer, &secondPlayer)

			if firstPlayer.output != tt.expectedFirstPlayerOutcome {
				t.Errorf("Error: expected outcome to be %v  but was %v", tt.expectedFirstPlayerOutcome, firstPlayer.output)
			}

			if secondPlayer.output != tt.expectedSecondPlayerOutcome {
				t.Errorf("Error: expected outcome to be %v  but was %v", tt.expectedFirstPlayerOutcome, firstPlayer.output)
			}
		})
	}
}
