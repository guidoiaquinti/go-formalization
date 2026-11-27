# Independent review of `internal/strconv.AppendBool`

**Verdict: FAITHFUL**

Evidence SHA-256: `866bc8cb0e3b5c778be2bd57b2187be1ed79f6558e3b0cd50a68abdb311b7e81`

I checked evidence commit `cf525f3` against Go 1.27.0 `src/internal/strconv/atob.go` and `src/internal/strconv/atoi.go`. Pinned Go appends the complete ASCII word as one operation. The fixtures distinguish exact-fit reuse, preservation of spare capacity bytes and detaching on growth without touching the original backing.

Pinned Go and Lean agreed on 4,416 fixed cases and 20,000 samples per seed using 17, 2718 and 65537 on darwin/arm64 with `go run ./difftest/internal/strconv/bool`. Five independent literal fixtures passed with `go run ./difftest/internal/strconv/bool/review` and detected ten wrong observations across the group. The scoped kernel audit passed. Scope covers arbitrary raw parse strings, both Boolean values, legal destination geometry and signed native error tags. Allocator-selected capacity after growth is outside scope. [Independent control](../../../../difftest/internal/strconv/bool/review/main.go)
