package strings

import "io"

func GoLeanReviewGetStringWriter(w io.Writer) io.StringWriter { return getStringWriter(w) }
