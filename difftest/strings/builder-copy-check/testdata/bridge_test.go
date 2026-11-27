package strings

func GoLeanBuilderCopyCheckSetup(mode, capacity int, s string) *Builder {
	if mode == 3 {
		return nil
	}
	b := new(Builder)
	if capacity >= 0 {
		b.buf = make([]byte, len(s), capacity)
		copy(b.buf, s)
	}
	switch mode {
	case 1:
		b.addr = b
	case 2:
		original := new(Builder)
		original.addr = original
		original.buf = b.buf
		*b = *original
	}
	return b
}

func GoLeanBuilderCopyCheck(b *Builder) { b.copyCheck() }

func GoLeanBuilderCopyCheckObserve(b *Builder) (identity string, addr *Builder, capacity int, s string) {
	identity = "unset"
	if b.addr == b {
		identity = "self"
	} else if b.addr != nil {
		identity = "foreign"
	}
	capacity = -1
	if b.buf != nil {
		capacity = cap(b.buf)
	}
	return identity, b.addr, capacity, string(b.buf)
}
