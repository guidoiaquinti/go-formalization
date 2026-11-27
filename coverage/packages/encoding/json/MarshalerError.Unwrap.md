# Independent review of `encoding/json.MarshalerError.Unwrap`

**Verdict: FAITHFUL**

Evidence SHA-256: `990edb63e3320015d45b2fd5fec099241d2d300f01acf170fc892f11a468003e`

I reviewed evidence commit `9c215c3` against the default Go 1.27.0 compatibility sources in `src/encoding/json/v2_scanner.go`, `v2_encode.go` and `v2_decode.go`. The method returns the stored error interface without invoking Error. Go checks identity for a live nested pointer and a typed nil. Lean carries an abstract identity token.

Pinned Go and Lean agreed on 1,009 fixed rows and 20,000 sampled rows per seed for 17, 2718 and 65537 on darwin/arm64 via `go run ./difftest/encoding/json/compat-errors`. Six independent literal fixtures passed with `go run ./difftest/encoding/json/compat-errors/review` and detected eleven incorrect outcomes across the group. The scoped kernel audit passed. Scope covers direct error methods and their default compatibility implementation. Reflection formatting internals and callback bodies are supplied inputs. [Independent control](../../../../difftest/encoding/json/compat-errors/review/main.go)
