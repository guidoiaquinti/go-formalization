# Independent review of `encoding/json.MarshalerError.Error`

**Verdict: FAITHFUL**

Evidence SHA-256: `9233d440dee0e5bc7a9bf6bb1eca2a151860cd2aba0631c45792577f017b6dd1`

I reviewed evidence commit `9c215c3` against the default Go 1.27.0 compatibility sources in `src/encoding/json/v2_scanner.go`, `v2_encode.go` and `v2_decode.go`. The method selects MarshalJSON when the stored source is empty and otherwise retains the source verbatim. It renders the type before invoking the nested Error method. The controls distinguish a nil type from a typed nil error that panics with the original token.

Pinned Go and Lean agreed on 1,009 fixed rows and 20,000 sampled rows per seed for 17, 2718 and 65537 on darwin/arm64 via `go run ./difftest/encoding/json/compat-errors`. Six independent literal fixtures passed with `go run ./difftest/encoding/json/compat-errors/review` and detected eleven incorrect outcomes across the group. The scoped kernel audit passed. Scope covers direct error methods and their default compatibility implementation. Reflection formatting internals and callback bodies are supplied inputs. [Independent control](../../../../difftest/encoding/json/compat-errors/review/main.go)
