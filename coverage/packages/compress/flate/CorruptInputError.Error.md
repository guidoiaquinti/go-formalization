# Independent review of `compress/flate.CorruptInputError.Error`

**Verdict: FAITHFUL**

The source converts the signed int64 offset to base ten and appends it to the exact corrupt-input prefix. The theorem fixes the negative-one case and controls include both signed endpoints.

The pinned Go 1.27.0 and Lean comparisons replayed 280 fixed cases and 60,000 sampled cases at seeds 17, 2718 and 65537 on darwin/arm64. The scoped Lean axiom audit passed. The model is scoped to signed offsets, raw strings and the measured nil and nested callback outcomes.

Nine independent controls exercise all four methods with pinned native Go and Lean. Using unsigned decimal formatting disagrees on 4 rows. Replacing malformed UTF-8 disagrees on 3 rows. Formatting the nested error through fmt, which turns its panic into text, disagrees on 1 row.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/compress/flate/error-text/review
```

Evidence SHA-256: `c91f09da3b64ee1234ac090778c55a30c7ec0f067c9327582d1efd705d1e3fa4`
