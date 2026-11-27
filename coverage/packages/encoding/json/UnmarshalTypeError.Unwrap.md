# Independent review of `encoding/json.UnmarshalTypeError.Unwrap`

**Verdict: FAITHFUL**

Evidence SHA-256: `d6875c605e702f1053ca0f56cd529ee2b173669d9a6a00623a59aca436bb2bbd`

I reviewed evidence commit `9c215c3` against the default Go 1.27.0 compatibility sources in `src/encoding/json/v2_scanner.go`, `v2_encode.go` and `v2_decode.go`. The method returns the exact immediate stored interface. The pinned native tests distinguish nil interface and typed nil without invoking the nested Error callback.

Pinned Go and Lean agreed on 1,009 fixed rows and 20,000 sampled rows per seed for 17, 2718 and 65537 on darwin/arm64 via `go run ./difftest/encoding/json/compat-errors`. Six independent literal fixtures passed with `go run ./difftest/encoding/json/compat-errors/review` and detected eleven incorrect outcomes across the group. The scoped kernel audit passed. Scope covers direct error methods and their default compatibility implementation. Reflection formatting internals and callback bodies are supplied inputs. [Independent control](../../../../difftest/encoding/json/compat-errors/review/main.go)
