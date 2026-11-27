# Independent review of `compress/flate.noEOF`

**Verdict: FAITHFUL**

Evidence SHA-256: `effb2d258326cf5207e9cfe42badf84fc7cbbf1873b3d68604ec3b28b11e508d`

I checked evidence commit `89dc043` against Go 1.27.0 `src/compress/flate/inflate.go`. Exact equality with the standard EOF sentinel selects the unexpected EOF sentinel. The second theorem proves that all other modeled identities pass through. Interface equality against the ordinary pointer-backed sentinel also leaves an error with an incomparable dynamic type unchanged.

I ran `go run ./difftest/compress/flate/no-eof` on darwin/arm64 and reproduced all eight pinned native and Lean cases. Five independent cases passed with `go run ./difftest/compress/flate/no-eof/review`. The control distinguishes joined EOF, an error with user hooks, a typed nil and an incomparable map-backed error whose backing is preserved. It also checks that native equality invokes none of the error hooks. Five deliberately wrong outcomes failed. The scope assumes the standard I/O sentinel values and covers result identity. Decompressor I/O and reassigned I/O sentinels are outside scope. [Independent control](../../../../difftest/compress/flate/no-eof/review/main.go)
