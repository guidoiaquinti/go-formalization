package utf8

// GoLeanWord calls the pinned internal helper for both generic instantiations.
func GoLeanWord[T string | []byte](s T) uintptr { return word(s) }
