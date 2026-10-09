package user

import "math"

// Target computation — authoritative implementation of erd.md §5 (Mifflin-St Jeor
// → TDEE → goal adjustment). Pure functions, no I/O, so they are unit-tested directly
// (spec TC-ON-U1). The daily calorie target is computed on-the-fly, never stored.

// activityMultipliers maps activity_level → TDEE multiplier (erd.md §5 step 2).
var activityMultipliers = map[string]float64{
	"sedentary":   1.2,
	"light":       1.375,
	"moderate":    1.55,
	"active":      1.725,
	"very_active": 1.9,
}

// goalAdjustments maps goal → kcal delta applied to TDEE (erd.md §5 step 3).
var goalAdjustments = map[string]float64{
	"lose":     -500,
	"maintain": 0,
	"gain":     300,
}

// TargetResult is the full computed breakdown (rounded to 2 decimals for the response).
type TargetResult struct {
	BMR                float64
	TDEE               float64
	DailyCalorieTarget float64
}

// BMR computes the Mifflin-St Jeor basal metabolic rate.
func BMR(sex string, weightKg, heightCm float64, age int) float64 {
	base := 10*weightKg + 6.25*heightCm - 5*float64(age)
	if sex == "female" {
		return base - 161
	}
	return base + 5 // male (and default)
}

// ComputeTarget runs the full BMR → TDEE → goal pipeline. Values are rounded to 2
// decimals. Unknown activity/goal (should be blocked by validation) fall back to
// neutral factors so the function never panics.
func ComputeTarget(sex string, weightKg, heightCm float64, age int, activityLevel, goal string) TargetResult {
	bmr := BMR(sex, weightKg, heightCm, age)

	mult, ok := activityMultipliers[activityLevel]
	if !ok {
		mult = 1.2
	}
	tdee := bmr * mult

	adj := goalAdjustments[goal] // zero value (0) is the safe maintain-like default
	target := tdee + adj

	return TargetResult{
		BMR:                round2(bmr),
		TDEE:               round2(tdee),
		DailyCalorieTarget: round2(target),
	}
}

func round2(v float64) float64 { return math.Round(v*100) / 100 }
