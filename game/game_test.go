package game

import (
	"fmt"
	"testing"
)

var playerTestValues = []map[string]Color{
	{},
	{
		"Ezekiel": Blue,
	},
	{
		"Susan":  Red,
		"Robert": Yellow,
	},
	{
		"Gabriel": Green,
		"Lori":    Blue,
		"Clinton": Yellow,
	},
	{
		"Sergio": Blue,
		"Traci":  Green,
		"Erik":   Red,
		"Emmett": Yellow,
	},
}

func TestGamePlayerCountDirect(t *testing.T) {
	testCases := map[int]int{
		0: 0,
		1: 1,
		2: 2,
		3: 3,
		4: 4,
	}

	for input, expected := range testCases {
		t.Run(fmt.Sprintf("%d players", input), func(t *testing.T) {
			g := Game{
				playerCount: input,
			}
			if actual := g.PlayerCount(); actual != expected {
				t.Fatalf("expected %d but got %d", expected, actual)
			}
		})
	}
}

func TestGamePlayersDirect(t *testing.T) {
	for i, inputs := range playerTestValues[1:] {
		count := i + 1
		t.Run(fmt.Sprintf("%d players", count), func(t *testing.T) {
			g := Game{
				playerCount: count,
			}
			for name, color := range inputs {
				player := Player{
					Name:  name,
					color: color,
				}
				g.players[i] = player
			}
			players := g.Players()
			expected, actual := len(inputs), len(players)
			if actual != expected {
				t.Fatalf("expected %d players but got %d players", expected, actual)
			}
			for _, player := range players {
				expected, actual := inputs[player.Name], player.color
				if actual != expected {
					t.Fatalf("expected player %s to have color %s but got color %s", player.Name, expected, actual)
				}
			}
		})
	}
}

func TestAvailableColors(t *testing.T) {
	testCases := []map[Color]bool{
		{
			Blue:   true,
			Green:  true,
			Red:    true,
			Yellow: true,
		},
		{
			Blue:   false,
			Green:  true,
			Red:    true,
			Yellow: true,
		},
		{
			Blue:   true,
			Green:  true,
			Red:    false,
			Yellow: false,
		},
		{
			Blue:   false,
			Green:  false,
			Red:    true,
			Yellow: false,
		},
		{
			Blue:   false,
			Green:  false,
			Red:    false,
			Yellow: false,
		},
	}

	for i, inputs := range playerTestValues {
		t.Run(fmt.Sprintf("%d players", i), func(t *testing.T) {
			g := Game{}
			for name, color := range inputs {
				g, _ = g.AddPlayer(name, color)
			}
			colors := g.AvailableColors()
			for color, availability := range testCases[i] {
				if expected, actual := availability, colors[color]; actual != expected {
					t.Fatalf("wrong availability for color %s; expected %t but got %t", color, expected, actual)
				}
			}
		})
	}
}

func TestGameAddPlayerValid(t *testing.T) {
	for i, inputs := range playerTestValues[1:] {
		count := i + 1
		t.Run(fmt.Sprintf("%d players", count), func(t *testing.T) {
			var (
				err error
				g   = Game{}
				j   int
			)

			for name, color := range inputs {
				if g, err = g.AddPlayer(name, color); err != nil {
					t.Fatalf("failed to add player: %s", err)
				}
				j++
				if g.playerCount != j {
					t.Fatalf("expected player count %d but got %d", j, g.playerCount)
				}
			}

			// check for player existence both ways
			players := g.Players()
			for playerIn := range inputs {
				var found bool
				for _, playerOut := range players {
					if playerOut.Name == playerIn {
						found = true
					}
				}
				if !found {
					t.Fatalf("player %s is missing", playerIn)
				}
			}
			for _, playerOut := range players {
				if inputs[playerOut.Name] == InvalidColor {
					t.Fatalf("player %s should not be in game", playerOut.Name)
				}
			}
		})
	}
}

func TestGameAddPlayerTooManyPlayers(t *testing.T) {
	var (
		err         error
		g           = Game{}
		playerColor = Blue
		playerName  = "Marilyn"
	)

	for name, color := range playerTestValues[4] {
		if g, err = g.AddPlayer(name, color); err != nil {
			t.Fatalf("failed to add player %s with color %s", name, color)
		}
	}
	if _, err = g.AddPlayer(playerName, playerColor); err == nil {
		t.Fatalf("shouldn't be able to add player %s with color %s", playerName, playerColor)
	}
}

func TestGameAddPlayerInvalidName(t *testing.T) {
	g := Game{}
	if _, err := g.AddPlayer("", Blue); err == nil {
		t.Fatalf("shouldn't be able to add player with empty name")
	}
}

func TestGameAddPlayerCollision(t *testing.T) {
	testCases := map[Color][]string{
		Blue: {
			"Rosa",
			"Noah",
		},
		Green: {
			"Marcus",
			"Billy",
		},
		Red: {
			"Terry",
			"Peggy",
		},
		Yellow: {
			"Josh",
			"Courtney",
		},
	}

	for color, names := range testCases {
		t.Run(fmt.Sprintf("color %s", color), func(t *testing.T) {
			var (
				err error
				g   = Game{}
			)

			for i := 0; i < len(names); i++ {
				name := names[i]
				g, err = g.AddPlayer(name, color)
				if i == 0 && err != nil {
					t.Fatalf("failed to add player %s with color %s", name, color)
				} else if i >= 1 && err == nil {
					t.Fatalf("shouldn't be able to add player %s with color %s", name, color)
				}
			}
		})
	}
}

func TestGameRemovePlayerValid(t *testing.T) {
	for i, inputs := range playerTestValues[1:] {
		count := i + 1
		t.Run(fmt.Sprintf("%d players", count), func(t *testing.T) {
			var (
				err error
				g   = Game{}
			)

			for name, color := range inputs {
				g, _ = g.AddPlayer(name, color)
			}
			for name := range inputs {
				if g, err = g.RemovePlayer(name); err != nil {
					t.Fatalf("failed to remove player %s", name)
				}
				count--
				if actual := g.PlayerCount(); actual != count {
					t.Fatalf("expected %d players after removing %s but got %d", count, name, actual)
				}
			}
		})
	}
}

func TestGameRemovePlayerNoPlayers(t *testing.T) {
	g := Game{}
	if _, err := g.RemovePlayer("Myra"); err == nil {
		t.Fatal("shouldn't be able to remove a player with no players")
	}
}

func TestGameRemovePlayerEmptyName(t *testing.T) {
	g := Game{}
	g, _ = g.AddPlayer("Fannie", Blue)
	if _, err := g.RemovePlayer(""); err == nil {
		t.Fatal("shouldn't be able to remove player with empty name")
	}
}

func TestGameRemovePlayerInvalidNames(t *testing.T) {
	testCases := [][]string{
		{
			"Grace",
		},
		{
			"Marvin",
			"Bobby",
		},
		{
			"Kayla",
			"Sally",
			"Ashley",
		},
		{
			"Keith",
			"Christie",
			"Lorenzo",
			"Rene",
		},
	}

	for i, inputs := range playerTestValues[1:] {
		t.Run(fmt.Sprintf("%d players", i+1), func(t *testing.T) {
			var (
				err error
				g   = Game{}
			)

			for name, color := range inputs {
				g, _ = g.AddPlayer(name, color)
			}
			for _, name := range testCases[i] {
				if g, err = g.RemovePlayer(name); err == nil {
					t.Fatalf("shouldn't be able to remove non-existant player %s", name)
				}
			}
		})
	}
}

func TestTakeTurn(t *testing.T) {
	var (
		testCases = map[string]Color{
			"Larisa Peyton":      Blue,
			"Ellsworth Delacruz": Red,
			"Merlyn Schulze":     Green,
			"Man Beckman":        Yellow,
			"Nguyet Samuels":     Green,
			"Emmitt Venable":     Red,
			"Sol Hite":           Green,
			"Ginny Allan":        Yellow,
			"Althea Cottrell":    Blue,
			"Florine Ragan":      Green,
			"Mckinley Bratcher":  Blue,
			"Hassan Cave":        Red,
			"Stefan Lind":        Yellow,
			"Faustino Martz":     Green,
			"Zenia Hutcheson":    Blue,
			"Denyse Roderick":    Yellow,
			"Anisha Downey":      Green,
			"Eufemia Merrick":    Blue,
			"Josue Goins":        Red,
			"Julian Groce":       Yellow,
			"Clarine Berlin":     Green,
			"Jong Carrillo":      Blue,
			"Jose Wild":          Red,
			"Debora Packard":     Yellow,
			"Taisha Locklear":    Green,
		}
		validSpins = make(map[int]bool)
	)

	for _, val := range spinnerValues {
		validSpins[val] = true
	}

	for name, color := range testCases {
		t.Run(fmt.Sprintf("%s %s", name, color), func(t *testing.T) {
			var (
				err    error
				player = Player{
					color: color,
					Name:  name,
				}
				turn Turn
			)

			switch turn, player, err = takeTurn(player); {
			case err != nil:
				t.Fatal("failed to generate a random number; local entropy is likely low")
			case !validSpins[turn.Spin]:
				t.Fatalf("turn has an invalid spin value: %d", turn.Spin)
			case turn.Player.Name != player.Name:
				t.Fatalf("player name should not change from %s to %s", player.Name, turn.Player.Name)
			case turn.Spin >= 0 && turn.Player.cherries != turn.Spin:
				t.Fatalf("expected %d cherries after turn but got %d", turn.Spin, turn.Player.cherries)
			case turn.Spin < 0 && turn.Player.cherries != 0:
				t.Fatalf("player should not have less than 0 cherries; got %d after spin %d", turn.Player.cherries, turn.Spin)
			}
		})
	}
}

func TestGamePlayNoPlayers(t *testing.T) {
	g := Game{}
	if _, _, err := g.Play(); err == nil {
		t.Fatal("shouldn't be able to play without any players")
	}
}

func TestGamePlayValid(t *testing.T) {
	for _, inputs := range playerTestValues[1:] {
		t.Run(fmt.Sprintf("%d players", len(inputs)), func(t *testing.T) {
			var (
				err    error
				g      = Game{}
				turns  []Turn
				winner Player
			)

			for name, color := range inputs {
				g, err = g.AddPlayer(name, color)
				if err != nil {
					t.Fatalf("failed to add player %s with color %s", name, color)
				}
			}
			if g.playerCount != len(inputs) {
				t.Fatalf("expected player count %d but got %d", len(inputs), g.playerCount)
			}

			switch turns, winner, err = g.Play(); {
			case turns == nil:
				t.Fatal("turn list is empty")
			case len(turns) < 2*g.playerCount+1:
				t.Fatalf("too few turns to have a winner; got %d turns", len(turns))
			case winner.cherries != 10:
				t.Fatalf("expected winner to have 10 cherries but got %d", winner.cherries)
			case err != nil:
				t.Fatalf("unexpected error: %s", err)
			}
		})
	}
}
