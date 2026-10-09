// Package unit holds fast, dependency-free unit tests (no HTTP, no DB).
// These cover pure logic that the contract (api-specs/user-service/01. Onboarding
// User.md) and ERD §5 define. Spec test cases: TC-ON-U1 (target) and TC-ON-U2 (hashing).
package unit

import (
	"math"
	"testing"

	"github.com/Monarqhs/cal-lowrisk/apps/backend/internal/modules/user"
	"github.com/stretchr/testify/assert"
)

// TC-ON-U1 — ComputeTarget matches ERD §5 across the full matrix of
// {both sexes} × {5 activity levels} × {3 goals}, using a fixed body so the
// expected numbers are hand-verifiable.
//
// Fixed body: 70 kg, 170 cm, 30 y.
//
//	male   BMR = 10*70 + 6.25*170 - 5*30 + 5   = 700 + 1062.5 - 150 + 5   = 1617.5
//	female BMR = 10*70 + 6.25*170 - 5*30 - 161 = 700 + 1062.5 - 150 - 161 = 1451.5
func TestComputeTarget_Matrix(t *testing.T) {
	const (
		weight = 70.0
		height = 170.0
		age    = 30
	)
	const maleBMR = 1617.5
	const femaleBMR = 1451.5

	multipliers := map[string]float64{
		"sedentary":   1.2,
		"light":       1.375,
		"moderate":    1.55,
		"active":      1.725,
		"very_active": 1.9,
	}
	goalDelta := map[string]float64{
		"lose":     -500,
		"maintain": 0,
		"gain":     300,
	}

	for _, sex := range []string{"male", "female"} {
		bmr := maleBMR
		if sex == "female" {
			bmr = femaleBMR
		}
		for activity, mult := range multipliers {
			for goal, delta := range goalDelta {
				sex, activity, goal := sex, activity, goal
				bmr, mult, delta := bmr, mult, delta
				name := sex + "/" + activity + "/" + goal
				t.Run(name, func(t *testing.T) {
					got := user.ComputeTarget(sex, weight, height, age, activity, goal)

					wantTDEE := round2(bmr * mult)
					wantTarget := round2(bmr*mult + delta)

					assert.Equal(t, round2(bmr), got.BMR, "BMR")
					assert.Equal(t, wantTDEE, got.TDEE, "TDEE")
					assert.Equal(t, wantTarget, got.DailyCalorieTarget, "DailyCalorieTarget")
				})
			}
		}
	}
}

// TC-ON-13 / spec worked example: male/72.5/175/29/moderate/lose -> target 2102.06.
func TestComputeTarget_WorkedExample(t *testing.T) {
	got := user.ComputeTarget("male", 72.5, 175, 29, "moderate", "lose")
	assert.Equal(t, 1678.75, got.BMR)
	assert.Equal(t, 2602.06, got.TDEE)
	assert.Equal(t, 2102.06, got.DailyCalorieTarget)
}

// TC-ON-14 — the female branch uses the -161 constant.
func TestBMR_FemaleBranch(t *testing.T) {
	male := user.BMR("male", 70, 170, 30)
	female := user.BMR("female", 70, 170, 30)
	assert.Equal(t, 1617.5, male)
	assert.Equal(t, 1451.5, female)
	assert.InDelta(t, 166.0, male-female, 0.0001, "male - female must be 5 - (-161) = 166")
}

// round2 mirrors the service's rounding (math.Round to 2 decimals) so expectations
// compare exactly.
func round2(v float64) float64 { return math.Round(v*100) / 100 }
