# Independent review of `internal/strconv.Error.Error`

**Verdict: FAITHFUL**

Evidence SHA-256: `538ca1eadaa105841aa8a17d8a093f4e9706e3412d3a7e6018b93027d93e863c`

I checked evidence commit `cf525f3` against Go 1.27.0 `src/internal/strconv/atob.go` and `src/internal/strconv/atoi.go`. Pinned Go maps tags one through four to their fixed diagnostic strings and all other native signed tags to unknown error. Independent cases cover each known tag and the native minimum.

Pinned Go and Lean agreed on 4,416 fixed cases and 20,000 samples per seed using 17, 2718 and 65537 on darwin/arm64 with `go run ./difftest/internal/strconv/bool`. Five independent literal fixtures passed with `go run ./difftest/internal/strconv/bool/review` and detected ten wrong observations across the group. The scoped kernel audit passed. Scope covers arbitrary raw parse strings, both Boolean values, legal destination geometry and signed native error tags. Allocator-selected capacity after growth is outside scope. [Independent control](../../../../difftest/internal/strconv/bool/review/main.go)
