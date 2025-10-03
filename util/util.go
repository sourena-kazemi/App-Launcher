package util

import (
	"math"
	"strings"
)

func CalculateHeight(itemCount int) float32 {
	inputHeight := float32(40)
	itemHeight := float32(40)

	if itemCount == 0 {
		return inputHeight
	}
	return inputHeight + (itemHeight * float32(math.Min(float64(itemCount), 10))) + (float32(math.Min(float64(itemCount), 10)+1) * 4)
}

func CleanExec(raw string) string {
	// Remove placeholders like %U, %f, etc.
	return strings.FieldsFunc(raw, func(r rune) bool {
		return r == '%' || r == '\n'
	})[0]
}
