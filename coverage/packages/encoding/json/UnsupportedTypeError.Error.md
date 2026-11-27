# Independent review of `encoding/json.UnsupportedTypeError.Error`

**Verdict: FAITHFUL**

Evidence SHA-256: `939fe988c160b6bad1eb8151c327baf42ff0f9e7e6ae9bbf0c3b26a738d53105`

I reviewed evidence commit `9c215c3` against the default Go 1.27.0 compatibility sources in `src/encoding/json/v2_scanner.go`, `v2_encode.go` and `v2_decode.go`. The method prefixes the reflected type name and panics when the type or receiver is nil. Lean treats the reflected type name as a supplied value.

Pinned Go and Lean agreed on 1,009 fixed rows and 20,000 sampled rows per seed for 17, 2718 and 65537 on darwin/arm64 via `go run ./difftest/encoding/json/compat-errors`. Six independent literal fixtures passed with `go run ./difftest/encoding/json/compat-errors/review` and detected eleven incorrect outcomes across the group. The scoped kernel audit passed. Scope covers direct error methods and their default compatibility implementation. Reflection formatting internals and callback bodies are supplied inputs. [Independent control](../../../../difftest/encoding/json/compat-errors/review/main.go)
