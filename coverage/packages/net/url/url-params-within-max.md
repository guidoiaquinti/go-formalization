# Independent review of `net/url.urlParamsWithinMax`

**Verdict: FAITHFUL**

Evidence SHA-256: `95d84988fc524ef81780dceba7b39e6c350fe5a33d2d2105728c08cf26ce1a2b`

Reviewed author commit `de11364` against pinned Go 1.27.0 `src/net/url/url.go:963–977`. The default comparison includes 10,000. A nonzero custom setting is exclusive, zero disables the limit, and an out-of-range signed override uses the default. The native control measures the metric increment around the exact call.

Eight fixed fixtures passed against pinned native Go and compiled Lean. Making the custom comparison inclusive changed one fixture. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/net/url/query/review` from the repository root.

Author evidence passed 2,463 boundaries and 60,000 sampled cases at seeds 17, 2718, and 65537. The model takes the effective GODEBUG value and native signed width as inputs. Process setting races, external backing alias changes, and allocation failure remain outside scope.
