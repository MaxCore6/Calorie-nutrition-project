package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
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
	WeightHistory    []WeightLog
	Activities       []Activity
}

type WeightLog struct {
	Date   string
	Weight float64
}

type Activity struct {
	Name         string
	Duration_min int
	Intensity    string
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
		activityFactor += 0.5
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

func weightChange(user AppUser) string {
	history := user.WeightHistory
	n := len(history)

	if n < 2 {
		return "Not enough data for analyze weight change"
	}

	last := history[n-1].Weight
	prev := history[n-2].Weight
	diff := last - prev

	if diff > 0 {
		return fmt.Sprintf("You have gained weight : %.1f kg since your recent weight", diff)
	} else if diff < 0 {
		return fmt.Sprintf("You have lost weight: %.1f kg since your recent weight", -diff)
	} else {
		return "Your weight has not changed since your last weight "
	}
}

func caloriesBurned(activity Activity, weight float64) float64 {
	metValues := map[string]float64{
		"Running": 9.8,
		"Walking": 3.5,
		"Swiming": 6.0,
		"Gym":     5.0,
	}

	met, ok := metValues[activity.Name]
	if !ok {
		met = 3.0
	}
	return met * weight * float64(activity.Duration_min) / 60.0
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
		WeightHistory: []WeightLog{
			{"2025-03-08", 99.0},
			{"2025-04-08", 97.3},
		},
		Activities: []Activity{
			{
				Name:         "Gym",
				Duration_min: 60,
				Intensity:    "high",
			},
		},
	}
	insertUser(user)

	calories := calculatingCaloriePerDay(user)
	fmt.Printf("Recommended daily calories for %s: %d kcal\n", user.Name, calories)
	fmt.Println(weightChange(user))

	for _, a := range user.Activities {
		burned := caloriesBurned(a, user.Weight)
		fmt.Printf("Activity: %s | Duration: %d min | Burned: %.2f kcal\n", a.Name, a.Duration_min, burned)
	}
}

func insertUser(user AppUser) {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/calories")
	if err != nil {
		fmt.Println("Unable to connect to database:", err)
		return
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(),
		"INSERT INTO users (name, age, weight, height, gender, activity, goal) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		"Max", 27, 98.5, 181.0, "male", "moderate", "lose_weight")
	if err != nil {
		fmt.Fprint(os.Stderr, "Insert failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("User inserted successfully!")
}
