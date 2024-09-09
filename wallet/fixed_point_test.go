package wallet

import (
	"math"
	"math/rand/v2"
	"testing"
)

// precisionTolerance is the maximum allowed difference between the original and
// converted float values
const precisionTolerance = 1e-10

// TestFixedPoint tests the fixed point implementation
func TestFixedPoint(t *testing.T) {
	// Generate a random float64 between 0 and 1000
	originalFloat := rand.Float64() * 1000

	// Convert to and from FixedPoint
	fixedPoint := FixedPointFromFloat(originalFloat)
	convertedFloat := fixedPoint.ToFloat()

	// Check if the converted value is within tolerance of the original
	if math.Abs(originalFloat-convertedFloat) > precisionTolerance {
		t.Errorf("Conversion not within tolerance. Original: %f, Converted: %f", originalFloat, convertedFloat)
	}
}
