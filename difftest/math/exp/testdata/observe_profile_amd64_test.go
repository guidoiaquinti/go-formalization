package math_test

import (
	"math"
	"os"
	"strconv"
	"testing"
)

func TestGoLeanExpProfile(t *testing.T) {
	if err := os.WriteFile(os.Getenv("GOLEAN_MATH_EXP_PROFILE"), []byte(strconv.FormatBool(math.GoLeanExpUsesFMA())), 0600); err != nil {
		t.Fatal(err)
	}
}
