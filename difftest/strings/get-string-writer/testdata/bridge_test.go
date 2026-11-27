package strings

import "io"

func GoLeanGetStringWriter(writer io.Writer) io.StringWriter { return getStringWriter(writer) }
