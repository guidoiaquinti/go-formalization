# Independent review of `encoding/json.UnsupportedValueError.Error`

**Verdict: FAITHFUL**

Evidence SHA-256: `2f367c8a8c11f7798074e59cc2e8fad019f7e826238cb7356d63d50df39e142d`

I reviewed evidence commit `9c215c3` against the default Go 1.27.0 compatibility sources in `src/encoding/json/v2_scanner.go`, `v2_encode.go` and `v2_decode.go`. The method prefixes the stored diagnostic. The control shows that the reflect.Value is ignored and that raw bytes survive.

Pinned Go and Lean agreed on 1,009 fixed rows and 20,000 sampled rows per seed for 17, 2718 and 65537 on darwin/arm64 via `go run ./difftest/encoding/json/compat-errors`. Six independent literal fixtures passed with `go run ./difftest/encoding/json/compat-errors/review` and detected eleven incorrect outcomes across the group. The scoped kernel audit passed. Scope covers direct error methods and their default compatibility implementation. Reflection formatting internals and callback bodies are supplied inputs. [Independent control](../../../../difftest/encoding/json/compat-errors/review/main.go)
