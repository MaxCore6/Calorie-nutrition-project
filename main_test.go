package main

import "testing"

func TestCalculatingCaloriePerDay(t *testing.T) {
	user := AppUser{
		Age:      27,
		Weight:   98.0,
		Height:   181.0,
		Gender:   "male",
		Activity: "high",
		Goal:     "maintance_weight",
	}

	expected := 98*10 + 6.25*181 - 5*27 + 5 // basic Bmr
	expected *= 1.75                        // activity
	got := calculatingCaloriePerDay(user)

	if got < int(expected)-50 || got > int(expected)+50 {
		t.Errorf("Expected around %d, got %d", int(expected), got)
	}
}

func TestGetActivityFactor(t *testing.T) {
	user := AppUser{
		Age:         27,
		Activity:    "high",
		StepsPerDay: 6700,
	}
}
