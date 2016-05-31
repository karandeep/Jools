package lib

import (
	"math"
)

const ELO_K = 24.0

// Calculate the expected % outcome
func Expected(Rb, Ra float64) float64 {
	return 1 / (1 + math.Pow(10, (Rb-Ra)/400))
}

// Calculate the new winner score
func Win(score, expected float64) float64 {
	return score + ELO_K*(1.0-expected)
}

// Calculate the new loser score
func Loss(score, expected float64) float64 {
	return score + ELO_K*(0.0-expected)
}
