package game

import (
	"fmt"
	"testing"
)

var testPlayers = []map[string]Color{
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

func TestGamePlayerCountIndirect(t *testing.T) {
	for expected, inputs := range testPlayers {
		t.Run(fmt.Sprintf("%d players", expected), func(t *testing.T) {
			g := Game{}
			for name, color := range inputs {
				g, _ = g.AddPlayer(name, color)
			}
			if actual := g.PlayerCount(); actual != expected {
				t.Fatalf("expected %d but got %d", expected, actual)
			}
		})
	}
}

func TestGamePlayersDirect(t *testing.T) {
	for i, inputs := range testPlayers[1:] {
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

func TestGamePlayersIndirect(t *testing.T) {
	for i, inputs := range testPlayers[1:] {
		count := i + 1
		t.Run(fmt.Sprintf("%d players", count), func(t *testing.T) {
			g := Game{}
			for name, color := range inputs {
				g, _ = g.AddPlayer(name, color)
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

func TestGameAddPlayerValid(t *testing.T) {
	for i, inputs := range testPlayers[1:] {
		count := i + 1
		t.Run(fmt.Sprintf("%d players", count), func(t *testing.T) {
			var (
				err error
				g   = Game{}
			)

			for name, color := range inputs {
				if g, err = g.AddPlayer(name, color); err != nil {
					t.Fatalf("failed to add player: %s", err)
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

	for name, color := range testPlayers[4] {
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
	for i, inputs := range testPlayers[1:] {
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

	for i, inputs := range testPlayers[1:] {
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
