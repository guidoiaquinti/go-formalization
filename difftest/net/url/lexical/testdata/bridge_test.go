package url

func GoLeanValidOptionalPort(s string) bool            { return validOptionalPort(s) }
func GoLeanSplitHostPort(s string) (string, string)    { return splitHostPort(s) }
func GoLeanGetScheme(s string) (string, string, error) { return getScheme(s) }
func GoLeanContainsCTL(s string) bool                  { return stringContainsCTLByte(s) }
