package strings_test

import (
	"strings"
	"testing"
)

func TestGoLeanBuildOnceControl(t *testing.T) {
	if failure := strings.GoLeanBuildOnceControl(); failure != "" {
		t.Fatal(failure)
	}
}
