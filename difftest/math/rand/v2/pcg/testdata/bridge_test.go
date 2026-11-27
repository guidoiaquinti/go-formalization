package rand

func GoLeanPCGState(p *PCG) (uint64, uint64) { return p.hi, p.lo }

func GoLeanPCGNext(p *PCG) (uint64, uint64) { return p.next() }

func GoLeanPCGUnmarshalError() error { return errUnmarshalPCG }
