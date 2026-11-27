# Independent review of `internal/strconv.FormatBool`

**Verdict: FAITHFUL**

Evidence SHA-256: `0d7a34e2bb7cf3b358c69747e68bc9e54450177a9a318a900d4c1b521bb96566`

I checked evidence commit `cf525f3` against Go 1.27.0 `src/internal/strconv/atob.go` and `src/internal/strconv/atoi.go`. Pinned Go selects the lowercase ASCII word for the Boolean. The Lean theorem lists both spellings, and the independent fixtures check each output.

Pinned Go and Lean agreed on 4,416 fixed cases and 20,000 samples per seed using 17, 2718 and 65537 on darwin/arm64 with `go run ./difftest/internal/strconv/bool`. Five independent literal fixtures passed with `go run ./difftest/internal/strconv/bool/review` and detected ten wrong observations across the group. The scoped kernel audit passed. Scope covers arbitrary raw parse strings, both Boolean values, legal destination geometry and signed native error tags. Allocator-selected capacity after growth is outside scope. [Independent control](../../../../difftest/internal/strconv/bool/review/main.go)
