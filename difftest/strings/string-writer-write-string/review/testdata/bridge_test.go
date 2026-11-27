package strings

import "io"

func GoLeanReviewStringWriter(w io.Writer, s string) (int, error) {
	return (stringWriter{w}).WriteString(s)
}
