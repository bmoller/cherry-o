package game

import (
	"fmt"
	"testing"
)

var (
	colorTestValues = map[Color]string{
		Blue:   "blue",
		Green:  "green",
		Red:    "red",
		Yellow: "yellow",
	}
)

func TestColorString(t *testing.T) {
	for color, name := range colorTestValues {
		t.Run(fmt.Sprint(name), func(t *testing.T) {
			if actual := color.String(); actual != name {
				t.Fatalf("expected %s but got %s", name, actual)
			}
		})
	}
}

func TestColorFilterValue(t *testing.T) {
	for color, name := range colorTestValues {
		t.Run(fmt.Sprint(name), func(t *testing.T) {
			if actual := color.FilterValue(); actual != name {
				t.Fatalf("expected %s but got %s", name, actual)
			}
		})
	}
}

func TestPlayerColor(t *testing.T) {
	testCases := map[string]Color{
		"Nga Barkley":    Blue,
		"Laticia Brower": Green,
		"Inell Mueller":  Red,
		"Cris Durr":      Yellow,
	}

	for name, color := range testCases {
		t.Run(fmt.Sprintf("%s %s", name, color), func(t *testing.T) {
			player := Player{
				color: color,
				Name:  name,
			}
			if actual := player.Color(); actual != color {
				t.Fatalf("expected %s but got %s", color, actual)
			}
		})
	}
}

func TestPlayerString(t *testing.T) {
	testCases := map[string]Color{
		"Catina Olivo":   Blue,
		"Elza Duff":      Green,
		"Delma Espinoza": Red,
		"Ardath Vang":    Yellow,
	}

	for name, color := range testCases {
		t.Run(fmt.Sprintf("%s %s", name, color), func(t *testing.T) {
			player := Player{
				color: color,
				Name:  name,
			}
			if actual := player.String(); actual != name {
				t.Fatalf("expected %s but got %s", name, actual)
			}
		})
	}
}

func TestPlayerFilterValue(t *testing.T) {
	testCases := map[string]Color{
		"Arleen Mello":      Blue,
		"Jeffie Nagel":      Green,
		"Vi Quinones":       Red,
		"Kiersten Stoddard": Yellow,
	}

	for name, color := range testCases {
		t.Run(fmt.Sprintf("%s %s", name, color), func(t *testing.T) {
			player := Player{
				color: color,
				Name:  name,
			}
			if actual := player.FilterValue(); actual != name {
				t.Fatalf("expected %s but got %s", name, actual)
			}
		})
	}
}

func TestPlayerUpdateCherries(t *testing.T) {
	testCases := []struct {
		name     string
		spins    []int
		expected int
	}{
		{
			name:     "Waylon Berlin",
			spins:    []int{-10},
			expected: 0,
		},
		{
			name:     "Earle Corbett",
			spins:    []int{-2, -2, 4, -10, 2, 1, 3},
			expected: 6,
		},
		{
			name:     "Gisele Wilkinson",
			spins:    []int{1, 1, 1, 1, 1, 1, 1, 1, 1},
			expected: 9,
		},
		{
			name:     "Yoshiko Weston",
			spins:    []int{-10, -10, -10, -10, -2, -2, 3},
			expected: 3,
		},
		{
			name:     "Shayne Dye",
			spins:    []int{4, 4, 1, 4},
			expected: 10,
		},
		{
			name:     "Tatyana Medrano",
			spins:    []int{2, -2, 2, -2, 2, -2, 2, -2},
			expected: 0,
		},
		{
			name:     "Darleen Lu",
			spins:    []int{3, 3, 3, -10, 3, 3, 3, -10, 3, 3, 3, -10, 1, 2, 3},
			expected: 6,
		},
		{
			name:     "Tawana Tennant",
			spins:    []int{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},
			expected: 10,
		},
		{
			name:     "Coretta Shearer",
			spins:    []int{3, 3, -2, -2, -2, -10, -2, -2, -2, -10, -2, -10, -10, 4},
			expected: 4,
		},
		{
			name:     "Missy Dickens",
			spins:    []int{2, 2, 2, 1, 1, 1, -10, 4, 4, -2, -2, -2},
			expected: 2,
		},
	}

	for _, test := range testCases {
		t.Run(fmt.Sprint(test.name), func(t *testing.T) {
			player := Player{
				Name: test.name,
			}
			for _, spin := range test.spins {
				player = player.updateCherries(spin)
			}
			if player.cherries != test.expected {
				t.Fatalf("expected %d but got %d", test.expected, player.cherries)
			}
		})
	}
}
