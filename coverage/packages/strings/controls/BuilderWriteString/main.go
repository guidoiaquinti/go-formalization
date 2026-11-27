package main

import (
	"fmt"
	"strings"
)

func main() {
	const copiedPanic = "strings: illegal use of non-zero Builder copied by value"
	var original strings.Builder
	zeroCopy := original
	if n, err := zeroCopy.WriteString("é\xff"); n != 3 || err != nil ||
		zeroCopy.String() != "é\xff" || original.String() != "" {
		panic("copy of untouched zero builder could not be initialized")
	}

	saved := zeroCopy.String()
	activeCopy := zeroCopy
	if failure := tryWrite(&activeCopy, "x"); failure != copiedPanic ||
		activeCopy.String() != saved || zeroCopy.String() != saved {
		panic("identified copy did not panic before append")
	}
	var empty strings.Builder
	if n, err := empty.WriteString(""); n != 0 || err != nil {
		panic("empty WriteString returned a nonzero count or an error")
	}
	emptyCopy := empty
	if failure := tryWrite(&emptyCopy, "x"); failure != copiedPanic {
		panic("empty write failed to initialize the copy-check identity")
	}
	if n, err := zeroCopy.WriteString("x"); n != 1 || err != nil ||
		zeroCopy.String() != saved+"x" || saved != "é\xff" {
		panic("ordinary append changed a previous string")
	}
	fmt.Println("Builder.WriteString: Go passed zero, live, copied and empty-write checks; copy-check omission rejected")
}

func tryWrite(b *strings.Builder, input string) (failure any) {
	defer func() { failure = recover() }()
	_, _ = b.WriteString(input)
	return nil
}
