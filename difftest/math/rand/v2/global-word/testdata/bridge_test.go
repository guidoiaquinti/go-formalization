package rand

// GoLeanWithGlobalSource installs a deterministic callback only for one
// serial overlay invocation and restores the actual package Source.
func GoLeanWithGlobalSource(source Source, observe func()) {
	saved := globalRand
	defer func() { globalRand = saved }()
	globalRand = New(source)
	observe()
}
