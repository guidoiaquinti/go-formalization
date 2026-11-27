# Independent review of `encoding/json.InvalidUnmarshalError.Error`

**Verdict: FAITHFUL**

Evidence SHA-256: `07523a23cb1ac3f0bffd3e022f74edeaa93dc9e38aa7b7ae37c291ba965efa77`

I reviewed evidence commit `9c215c3` against the default Go 1.27.0 compatibility sources in `src/encoding/json/v2_scanner.go`, `v2_encode.go` and `v2_decode.go`. The method uses three distinct branches for nil reflected type, pointer kind and other reflected kinds. The control also distinguishes a nil reflected type from a nil receiver.

Pinned Go and Lean agreed on 1,009 fixed rows and 20,000 sampled rows per seed for 17, 2718 and 65537 on darwin/arm64 via `go run ./difftest/encoding/json/compat-errors`. Six independent literal fixtures passed with `go run ./difftest/encoding/json/compat-errors/review` and detected eleven incorrect outcomes across the group. The scoped kernel audit passed. Scope covers direct error methods and their default compatibility implementation. Reflection formatting internals and callback bodies are supplied inputs. [Independent control](../../../../difftest/encoding/json/compat-errors/review/main.go)
