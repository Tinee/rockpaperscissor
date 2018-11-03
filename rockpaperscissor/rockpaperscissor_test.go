package rockpaperscissor_test

// import (
// 	"rockpaperscissor/rockpaperscissor"
// 	"testing"
// )

// type MockPlayer struct {
// 	name string
// 	move rockpaperscissor.Move
// }

// func (m MockPlayer) Move() rockpaperscissor.Move { return m.move }
// func (m MockPlayer) Name() string                { return m.name }

// func TestWeapon_Beats(t *testing.T) {
// 	const admin = "Admin Name"
// 	const opponent = "Opponent Name"

// 	tests := []struct {
// 		name           string
// 		admin          MockPlayer
// 		opponent       MockPlayer
// 		expectedWinner string
// 		draw           bool
// 	}{
// 		{
// 			"Paper beats rock",
// 			MockPlayer{
// 				admin,
// 				rockpaperscissor.Paper,
// 			},
// 			MockPlayer{
// 				opponent,
// 				rockpaperscissor.Rock,
// 			},
// 			admin,
// 			false,
// 		},
// 		{
// 			"Rock beats scissors",
// 			MockPlayer{
// 				admin,
// 				rockpaperscissor.Rock,
// 			},
// 			MockPlayer{
// 				opponent,
// 				rockpaperscissor.Scissor,
// 			},
// 			admin,
// 			false,
// 		},
// 		{
// 			"Scissor beats paper",
// 			MockPlayer{
// 				admin,
// 				rockpaperscissor.Scissor,
// 			},
// 			MockPlayer{
// 				opponent,
// 				rockpaperscissor.Paper,
// 			},
// 			admin,
// 			false,
// 		},
// 		{
// 			"Paper lose against scissors",
// 			MockPlayer{
// 				admin,
// 				rockpaperscissor.Paper,
// 			},
// 			MockPlayer{
// 				opponent,
// 				rockpaperscissor.Scissor,
// 			},
// 			opponent,
// 			false,
// 		},
// 		{
// 			"Rock lose against paper",
// 			MockPlayer{
// 				admin,
// 				rockpaperscissor.Rock,
// 			},
// 			MockPlayer{
// 				opponent,
// 				rockpaperscissor.Paper,
// 			},
// 			opponent,
// 			false,
// 		},
// 		{
// 			"Scissor lose against rock",
// 			MockPlayer{
// 				admin,
// 				rockpaperscissor.Scissor,
// 			},
// 			MockPlayer{
// 				opponent,
// 				rockpaperscissor.Rock,
// 			},

// 			opponent,
// 			false,
// 		},
// 		{
// 			"Paper draws against paper",
// 			MockPlayer{
// 				admin,
// 				rockpaperscissor.Paper,
// 			},
// 			MockPlayer{
// 				opponent,
// 				rockpaperscissor.Paper,
// 			},
// 			"",
// 			true,
// 		},
// 		{
// 			"Rock draws against rock",
// 			MockPlayer{
// 				admin,
// 				rockpaperscissor.Rock,
// 			},
// 			MockPlayer{
// 				opponent,
// 				rockpaperscissor.Rock,
// 			},
// 			"",
// 			true,
// 		},
// 		{
// 			"Scissor draws against scissor",
// 			MockPlayer{
// 				admin,
// 				rockpaperscissor.Scissor,
// 			},
// 			MockPlayer{
// 				opponent,
// 				rockpaperscissor.Scissor,
// 			},
// 			"",
// 			true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			winner, _, draw := rockpaperscissor.Play(tt.admin, tt.opponent)

// 			if tt.draw && !draw {
// 				t.Errorf("Error: expected it to be a draw")
// 			}

// 			if tt.draw {
// 				return
// 			}
// 			if name := winner.(MockPlayer).Name(); name != tt.expectedWinner {
// 				t.Errorf("Error: expected %v to win but %v won..", tt.expectedWinner, name)
// 			}
// 		})
// 	}
// }
