# Independent review of `net/url.ParseQuery`

**Verdict: FAITHFUL**

Evidence SHA-256: `01a817188aca2d26995eb293c83feeb48cea00f71a5a78ef50c2e6fd6c89c99a`

Reviewed author commit `de11364` against pinned Go 1.27.0 `src/net/url/url.go:952–956`. Initializes a nonnil map and returns partial data with the parser error. Fixed fixtures cover empty input, limit rejection, exact malformed-key payload, decoded semicolon in a key, and repeated values.

Eight fixed fixtures passed against pinned native Go and compiled Lean. Starting with a nil map changed all eight fixtures across the shared public observers. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/net/url/query/review` from the repository root.

Author evidence passed 2,463 boundaries and 60,000 sampled cases at seeds 17, 2718, and 65537. The model takes the effective GODEBUG value and native signed width as inputs. Process setting races, external backing alias changes, and allocation failure remain outside scope.
