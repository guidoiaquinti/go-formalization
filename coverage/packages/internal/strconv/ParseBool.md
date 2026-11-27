# Independent review of `internal/strconv.ParseBool`

**Verdict: FAITHFUL**

Evidence SHA-256: `d74f0c957949054b7b2c8a28dbc162c0e52b704a49127830e4332104daeb76d0`

I checked evidence commit `cf525f3` against Go 1.27.0 `src/internal/strconv/atob.go` and `src/internal/strconv/atoi.go`. Pinned Go accepts exactly six spellings of true and six of false. The native overlay checks bare ErrSyntax identity for unrecognized bytes. The fixed controls distinguish an accepted False, an unaccepted mixed case, whitespace and malformed UTF-8. Lean carries error presence.

Pinned Go and Lean agreed on 4,416 fixed cases and 20,000 samples per seed using 17, 2718 and 65537 on darwin/arm64 with `go run ./difftest/internal/strconv/bool`. Five independent literal fixtures passed with `go run ./difftest/internal/strconv/bool/review` and detected ten wrong observations across the group. The scoped kernel audit passed. Scope covers arbitrary raw parse strings, both Boolean values, legal destination geometry and signed native error tags. Allocator-selected capacity after growth is outside scope. [Independent control](../../../../difftest/internal/strconv/bool/review/main.go)
