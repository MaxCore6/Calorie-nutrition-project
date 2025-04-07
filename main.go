package main

import (
	"fmt"
)

type AppUser struct {
	Name             string  // name and surname
	Age              int     // age
	Weight           float64 // kg
	Height           float64 // cm
	Gender           string  // male / female
	Activity         string  // "low", "medium", "high"
	Goal             string  // "fast_lose_weight", "normal_lose_weight", "maintance_weight", "gain_weight"
	StepsPerDay      int     // quantity
	CardioPerWeek    int     // in minutes
	StrengthTraining int     // in minutes
	LevelInSports    string  // Amateur / Professional /
}

func getActivityFactor(user AppUser) float64 {
	activityLevels := map[string]float64{
		"low":    1.2,  // sedentary lifestyle
		"medium": 1.55, // average activity (3-4 workout per week)
		"high":   1.75, // high activity (5-7 workout per week)
	}
	activityFactor, exists := activityLevels[user.Activity]
	if !exists {
		activityFactor = 1.2 // default if not selected
	}

	if user.StepsPerDay == 0 && user.CardioPerWeek == 0 && user.StrengthTraining == 0 {
		return activityFactor // If all parameters 0, we return base coefficient
	}

	if user.StepsPerDay < 5000 {
		activityFactor -= 0.1
	} else if user.StepsPerDay > 10000 {
		activityFactor += 0.1
	}

	if user.CardioPerWeek < 60 {
		activityFactor -= 0.1
	} else if user.CardioPerWeek > 180 {
		activityFactor += 0.1
	}

	if user.StrengthTraining < 60 {
		activityFactor -= 0.1
	} else if user.StrengthTraining > 180 {
		activityFactor += 0.1
	}

	if activityFactor < 1.1 {
		activityFactor = 1.1
	} else if activityFactor > 2.0 {
		activityFactor = 2.0
	}

	return activityFactor
}

func calculatingCaloriePerDay(user AppUser) int {
	var bmr float64

	// calculation of bmr
	if user.Gender == "male" {
		bmr = 10*user.Weight + 6.25*user.Height - 5*float64(user.Age) + 5
	} else {
		bmr = 10*user.Weight + 6.25*user.Height - 5*float64(user.Age) - 161
	}

	bmr *= getActivityFactor(user)

	switch user.Goal {
	case "fast_lose_weight":
		bmr *= 0.75
	case "normal_lose_weight":
		bmr *= 0.85
	case "gain_weight":
		bmr *= 1.15
	}
	return int(bmr)
}

func main() {
	user := AppUser{
		Name:             "Maksim Makarov",
		Age:              27,
		Weight:           98.0,
		Height:           181.0,
		Gender:           "male",
		Activity:         "medium",
		Goal:             "normal_lose_weight",
		StepsPerDay:      6733,
		CardioPerWeek:    120,
		StrengthTraining: 360,
		LevelInSports:    "Amateur",
	}

	calories := calculatingCaloriePerDay(user)
	fmt.Printf("Recommended daily calories for %s: %d kcal\n", user.Name, calories)
}
